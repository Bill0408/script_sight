import torch
from .model import Model

# Use cuda if a gpu supporting cuda is available, if it isn't available, use the cpu.
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

model = Model()

model_path = "/script_sight/ai/model_weights.pth"
model.load_state_dict(torch.load(model_path, map_location=device))
model.to(device).eval()  # Move to GPU if available and set to evaluation mode
