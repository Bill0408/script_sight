from django.http import HttpResponse
from .forms import ImageUploadForm


#  predict handles incoming requests for predicting the right number in an image.
def predict(request):

    # Check if the request method is POST.
    if request.method == "POST":  # Only POST requests are allowed.

        # Create an instance of ImageUploadForm with the posted data and files.
        form = ImageUploadForm(request.POST, request.FILES)

        if form.is_valid():

            # Extract the uploaded image from the cleaned data of the form.
            modified_img = form.cleaned_data["uploadFile"]

            # Print the modified image.
            print(modified_img)

            # Send a successful HTTP response with an empty content.
            return HttpResponse("", content_type="text/plain")
        else:
            # If the form is not valid, print the error and send an error response.
            print(f'Invalid Image: {form.errors}')
            return HttpResponse(f'Invalid Image: {form.errors}', status=400, content_type="text/plain")

    # If the request method is not POST, send a 'Method Not Allowed' response.
    else:
        return HttpResponse('Method Not Allowed', status=405, content_type="text/plain")
