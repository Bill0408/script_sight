let cnv1 = document.getElementById("cnv1");
let cnv2 = document.getElementById("cnv2");
let predictionButton = document.getElementById("prediction");
let isDrawing = false;

// Get a 2d context from canvases because the drawings are all in 2d.
let ctx1 = cnv1.getContext("2d")
let ctx2 = cnv2.getContext("2d");

ctx2.font = "100px Arial"

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

    ctx1.lineWidth = 14;
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

predictionButton.addEventListener("click", async function () {
    if (predictionButton.innerHTML === "Make Prediction") {
        // Get the offscreen canvas and its context.
        let cnvCtx = createOffScreenCanvas();
        // Convert the image in the offscreen canvas to grayscale.

        // Get a data url
        let imageURL = cnvCtx.cnv.toDataURL("image/png");

        const responseData = await uploadImageUrl(imageURL);

        // Prepare the canvas
        ctx2.fillStyle = 'black';
        ctx2.fillRect(0, 0, cnv2.width, cnv2.height);

        // Increase font size and set style
        ctx2.font = "100px Arial";

        // Set the fill color for the text
        ctx2.fillStyle = 'white';

        // Center the number on the canvas
        ctx2.textBaseline = "middle"
        ctx2.textAlign = "center"

        ctx2.fillText(responseData, cnv2.width / 2, cnv2.height / 2);

        predictionButton.innerHTML = "Clear"

        // Draw the text
        ctx2.fillText(responseData, x, y);

        predictionButton.innerHTML = "Clear"
    } else {
        predictionButton.innerHTML = "Make Prediction"
        ctx1.fillStyle = 'black';
        ctx1.fillRect(0, 0, cnv1.width, cnv1.height);
        ctx2.fillStyle = 'black';
        ctx2.fillRect(0, 0, cnv2.width, cnv2.height);
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

// Upload a json containing the data url of the image to the go server.
async function uploadImageUrl(dataUrl) {
    let data = {
        imgUrl: dataUrl
    };
    try {
        const response = await fetch("upload", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            console.error("Response not OK:", response);
            return null;
        }

      return await response.text()

    } catch (error) {
        console.error("An error occurred:", error);
        return null;
    }
}
