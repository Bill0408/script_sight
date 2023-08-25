FROM python:3.9

# Set environment variables
# Ensures Python outputs everything that's printed inside the application rather than buffering it.
ENV PYTHONUNBUFFERED 1
# Set the default settings module for the Django application
ENV DJANGO_SETTINGS_MODULE ai_server.settings

# Set the working directory inside the container
WORKDIR /script_sight

# Copy the Django application to the container
COPY /ai_server ./

# Install Python dependencies
RUN pip install -r requirements.txt

# Expose port 8000 for the Django application
EXPOSE 8000

# Command to run the Django development server
CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
