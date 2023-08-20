let cnv1 = document.getElementById("cnv1");
let cnv2 = document.getElementById("cnv2");
let predictionButton = document.getElementById("prediction");
let isDrawing = false;

// Get a 2d context from canvases because the drawings are all in 2d.
let ctx1 = cnv1.getContext("2d")
let ctx2 = cnv2.getContext("2d");

// Fill the drawing canvas with a black rectangle that serves as the background color.
ctx1.fillStyle = 'black';
ctx1.fillRect(0, 0, cnv1.width, cnv1.height);

// getMousePos calculates the mouse position relative to the canvas.
function getMousePos(canvas, event) {
    // Get the rectangle that contains the canvas, it's margin, padding, etc in pixels.
    let rect = canvas.getBoundingClientRect();
    let scaleX = canvas.width / rect.width; // Get the scale of the canvas and rectangle's width.
    let scaleY = canvas.height / rect.height; // Get the scale of the canvas and rectangle's height.

    // Get the mouse position, subtract it from the rectangle position and multiply
    // the difference by the scale to get the mouse position relative to the canvas.
    let x = (event.clientX - rect.left) * scaleX;
    let y = (event.clientY - rect.top) * scaleY;

    return {
        x: x,
        y: y
    };
}

// When the mouse is pressed on the canvas, this begins creating
// new paths, and points at the mouse position relative to the canvas.
cnv1.addEventListener("mousedown", function(event) {
    isDrawing = true;
    let pos = getMousePos(cnv1, event);

    ctx1.beginPath();
    ctx1.moveTo(pos.x, pos.y);
});

// This uses the mouse position relative to the canvas to
// draw really short lines that look continuous as the mouse is moved.
function draw(event) {
    if (!isDrawing) return;

    let pos = getMousePos(cnv1, event);

    ctx1.lineWidth = 8;
    ctx1.lineCap = "round";
    ctx1.strokeStyle = "white";

    ctx1.lineTo(pos.x, pos.y);
    ctx1.stroke();
}

// Stop drawing the mouse left button is released.
cnv1.addEventListener("mouseup", function() {
    isDrawing = false
})

// As the mouse moves, call the draw function to start drawing.
cnv1.addEventListener('mousemove', draw)

predictionButton.addEventListener("click", function () {
    if (predictionButton.innerHTML === "Make Prediction") {
        // Get the offscreen canvas and its context.
        let cnvCtx = createOffScreenCanvas();
        // Convert the image in the offscreen canvas to grayscale.
        convertToGrayscale(cnvCtx.ctx);

        // Get a data url
        let imageURL = cnvCtx.cnv.toDataURL("image/png");

        uploadImageUrl(imageURL)

        predictionButton.innerHTML = "Clear"
    } else {
        predictionButton.innerHTML = "Make Prediction"
        ctx1.fillStyle = 'black';
        ctx1.fillRect(0, 0, cnv1.width, cnv1.height);
    }
});

function createOffScreenCanvas() {
    // Create an offscreen canvas for the resizing and grayscale conversion.
    let offscreenCanvas = document.createElement('canvas');
    offscreenCanvas.width = 28;
    offscreenCanvas.height = 28;
    let offscreenCtx = offscreenCanvas.getContext('2d');

    // Resize the original image
    offscreenCtx.drawImage(cnv1, 0, 0, 28, 28);

    return {
        cnv: offscreenCanvas,
        ctx: offscreenCtx
    }
}

function convertToGrayscale(ctx) {
    // Get image data for the specified region.
    let imageData = ctx.getImageData(0, 0, 28, 28);
    let data = imageData.data;

    // Iterate over pixel values in sets of 4 (R, G, B, A).
    for (let i = 0; i < data.length; i += 4) {
        // Calculate grayscale value using luminosity formula.
        let grayscale = 0.299 * data[i] + 0.587 * data[i + 1] + 0.114 * data[i + 2];

        // Set RGB channels to grayscale value.
        data[i] = grayscale;     // Red
        data[i + 1] = grayscale; // Green
        data[i + 2] = grayscale; // Blue
    }

    // Update the canvas with grayscale image data.
    ctx.putImageData(imageData, 0, 0);
}

async function uploadImageUrl(dataUrl) {
    let data = {
        imgUrl: dataUrl
    }

    try {
        const response = await fetch("upload", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            console.log(response)
        }

        const responseData = await response.json();
        console.log(responseData);
    } catch (error) {

    }
}