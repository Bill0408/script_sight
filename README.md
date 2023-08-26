# script_sight

## Table of Contents

- Introduction
- How it works
- Installation
- Demo
- Technologies Used

## Introduction

Script Sight is an interactive project where you can draw a digit on an HTML canvas, and an AI model will predict what digit you drew. 
It consists of two HTML canvases; one for drawing and the other for displaying the AI's predicted digit. 
The application is backed by a Go server that handles requests and a Django server running a PyTorch model for digit recognition.

## How it works

When you draw on the left canvas and click the "Make Prediction" button, an offscreen 28x28 HTML canvas is created to resize your drawing. This specific 28x28 dimension matches the input on which the AI, trained using the MNIST dataset, expects. This resizing ensures consistency and accuracy in prediction. The resized image is converted to a Data URL, capturing the PNG image in a base64-encoded format. This is sent to the Go server, which decodes it into an actual PNG image. This image is then sent in a multipart form request to the Django server running the AI. Before the prediction, the image undergoes several transformations: it's converted to grayscale for computational simplicity and to mimic the MNIST training data, transformed into a tensor (the standard format for AI processing), and its pixel values are normalized to ensure consistent scaling with the training data. The AI then predicts the digit, relays the information back to the Go server, which finally sends the prediction to be displayed on the right canvas in your browser.

## Installation

### Prerequisites

You only need Docker installed to run this project. I used docker to containerize the whole thing, 
so that you don't have to download the necessary languages, frameworks and dependencies. You can
download docker at this link: https://www.docker.com/products/docker-desktop/

### Steps

#### 1. Clone this repository

Open your terminal (or Command Prompt on Windows) and run the following command to clone the repository:

```bash
https://github.com/Bill0408/script_sight.git
```

#### 2. Navigate into the project directory

Type the following comand to navigate to the project directory:

```bash
cd script-sight
```

#### 3. Run Docker Compose to start the services

Type the following command to run Docker Compose:

```bash
docker-compose up
```

#### 4. Open your web browser and then copy and paste this: `http://localhost:8080`

## Demo

![Demo](/demo.gif)

## Technologies Used

### Languages
- Go
- JavaScript
- Python
- Dockerfile
- HTML
### Frameworks
- Bootstrap
- PyTorch
- Django
