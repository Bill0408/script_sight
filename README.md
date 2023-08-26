# script_sight

## Table of Contents

- Introduction
- Installation
- Demo
- Technologies Used

## Introduction

Script Sight is an interactive project where you can draw a digit on an HTML canvas, and an AI model will predict what digit you drew. 
It consists of two HTML canvases; one for drawing and the other for displaying the AI's predicted digit. 
The application is backed by a Go server that handles requests and a Django server running a PyTorch model for digit recognition.

## Installation

### Prerequisites

You only need Docker installed to run this project. I used docker to containerize the whole thing, 
so that you don't have to download the necessary languages, frameworks and dependencies. You can
download docker at this link: https://www.docker.com/products/docker-desktop/

### Steps

#### 1. Clone this repository

Open your terminal (or Command Prompt on Windows) and run the following command to clone the repository:

```bash
git clone https://github.com/yourusername/script-sight.git
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
