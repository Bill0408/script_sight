import torch
from django.http import HttpResponse
from torchvision.transforms import transforms
from PIL import Image  # Make sure to import this
from .forms import ImageUploadForm
from .model_loading import model

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")


def predict(request):
    try:
        if request.method == "POST":

            form = ImageUploadForm(request.POST, request.FILES)

            if form.is_valid():
                # Retrieve the uploaded image from the multipart form.
                uploaded_file = form.cleaned_data["uploadFile"]

                # Convert InMemoryUploadedFile to a PIL Image
                modified_img = Image.open(uploaded_file)

                # Add some transformations to the image: first, convert the image from RGB, which has 3
                # color channels. to grayscale, which has 1 color channel. This is done because the model
                # was trained on grayscale images. Next, Convert the image to a tensor because it is
                # the datatype that the neural network works with. Finally, apply, normalize the
                # pixel values of the tensor, so that they are in the range -1 to 1. This is useful
                # as values from -1 to 1, and 0 as the mean, requires less computation.
                transform = transforms.Compose([
                    transforms.Grayscale(num_output_channels=1),
                    transforms.ToTensor(),
                    transforms.Normalize((0.5,), (0.5,))
                ])

                # Apply the transformations to the image, and then add a new dimension at the beginning,
                # so that it is a batch. This is done because the neural network expect batches of images.
                image_tensor = transform(modified_img).unsqueeze(0).to(device)

                # Run the model with no gradient calculation because
                # it is not being trained. It is only making prediction.
                with torch.no_grad():
                    outputs = model(image_tensor)
                    _, predicted = outputs.max(1)

                # Send back the number that the model predicted.
                return HttpResponse(predicted.item(), content_type="text/plain")
            else:
                print(f'Invalid Image: {form.errors}')
                return HttpResponse(f'Invalid Image: {form.errors}', status=400, content_type="text/plain")
        else:
            return HttpResponse('Method Not Allowed', status=405, content_type="text/plain")

    except Exception as e:
        print(f"An error occurred: {e}")
        return HttpResponse(f"An error occurred: {e}", status=500, content_type="text/plain")
