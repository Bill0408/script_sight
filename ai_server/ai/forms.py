from django import forms


# ImageUploadForm represents a Django form used for uploading images.
# It inherits from Django's Form class.
class ImageUploadForm(forms.Form):
    # Declaring an ImageField named 'uploadFile' to handle image uploads.
    uploadFile = forms.ImageField()
