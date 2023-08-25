import torch
from torchvision import datasets, transforms
import torch.nn as nn

# Use cuda if a gpu supporting cuda is available, if it isn't available, use the cpu.
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# Preprocessing: Tensor conversion and Normalization. Preps data for better training performance.
transform = transforms.Compose([transforms.ToTensor(), transforms.Normalize((0.5,), (0.5,))])

# MNIST dataset for training. Automatically downloads and applies transformations.
train_dataset = datasets.MNIST(root='./data', train=True, download=True, transform=transform)
# DataLoader to handle batching of the training dataset.
train_loader = torch.utils.data.DataLoader(train_dataset, batch_size=1000, shuffle=True)

# MNIST dataset for testing. Automatically downloads and applies transformations.
test_dataset = datasets.MNIST(root='./data', train=False, download=True, transform=transform)
# DataLoader to handle batching of the test dataset.
test_loader = torch.utils.data.DataLoader(test_dataset, batch_size=1000, shuffle=False)


class Model(nn.Module):
    def __init__(self):
        super().__init__()

        # Adding multiple Conv-BatchNorm blocks to extract more complex features from the input.
        self.conv1 = nn.Conv2d(1, 64, kernel_size=3, padding=1)
        self.bn1 = nn.BatchNorm2d(64)
        self.conv2 = nn.Conv2d(64, 128, kernel_size=3, padding=1)
        self.bn2 = nn.BatchNorm2d(128)
        self.conv3 = nn.Conv2d(128, 256, kernel_size=3, padding=1)
        self.bn3 = nn.BatchNorm2d(256)

        # ReLU is used for the activation function because it helps with gradient propagation.
        self.relu = nn.ReLU()

        # Max-pooling is used to reduce spatial dimensions (width & height), making the network less computationally
        # intensive.
        self.max_pool = nn.MaxPool2d(2, stride=2)

        # Dropout is added for regularization, to prevent overfitting.
        self.drop = nn.Dropout(0.5)

        # Fully connected layers to perform classification based on the features extracted by the convolution layers.
        self.fc1 = nn.Linear(256 * 3 * 3, 128)
        self.fc2 = nn.Linear(128, 10)

    def forward(self, x):
        # Pass through the first Conv-BatchNorm-Activation-Pooling block; enhances basic features like edges
        x = self.max_pool(self.relu(self.bn1(self.conv1(x))))
        x = self.drop(x)

        # Pass through the second Conv-BatchNorm-Activation-Pooling block; enhances complex features
        x = self.max_pool(self.relu(self.bn2(self.conv2(x))))
        x = self.drop(x)

        # Pass through the third Conv-BatchNorm-Activation-Pooling block; further refines complex features
        x = self.max_pool(self.relu(self.bn3(self.conv3(x))))
        x = self.drop(x)

        # Flatten the tensor for Fully Connected layer input.
        x = x.view(x.size(0), -1)

        # Pass through the first Fully Connected layer and activate; compresses spatial features to feature vector
        x = self.relu(self.fc1(x))
        x = self.drop(x)

        # Pass through the final Fully Connected layer; maps feature vectors to one of 10 classes
        x = self.fc2(x)

        return x


# Training loop: processes the dataset in batches and updates model weights.
def train_loop(dataloader, model, loss_fn, optimizer):
    size = len(dataloader.dataset)
    model.train()

    for batch, (X, y) in enumerate(dataloader):
        X, y = X.to(device), y.to(device)
        pred = model(X)
        loss = loss_fn(pred, y)

        # Backpropagation and weight update steps.
        loss.backward()
        optimizer.step()
        optimizer.zero_grad()

        # Log training progress.
        if batch % 100 == 0:
            loss, current = loss.item(), (batch + 1) * len(X)
            print(f"loss: {loss:>7f}  [{current:>5d}/{size:>5d}]")


# Test loop: Evaluates model performance on a separate dataset not seen during training.
def test_loop(dataloader, model, loss_fn):
    model.eval()
    size = len(dataloader.dataset)
    num_batches = len(dataloader)
    test_loss, correct = 0, 0

    with torch.no_grad():
        for X, y in dataloader:
            X, y = X.to(device), y.to(device)
            pred = model(X)
            test_loss += loss_fn(pred, y).item()
            correct += (pred.argmax(1) == y).type(torch.float).sum().item()

    # Calculate and display test accuracy and average loss.
    test_loss /= num_batches
    correct /= size
    print(f"Test Error: \n Accuracy: {(100 * correct):>0.1f}%, Avg loss: {test_loss:>8f} \n")


# Initialize the model, loss, and optimizer.
model = Model().to(device)
loss = nn.CrossEntropyLoss()
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)

# Loop over epochs, invoking the training and test loops.
epochs = 10
for t in range(epochs):
    print(f"Epoch {t + 1}\n-------------------------------")
    train_loop(train_loader, model, loss, optimizer)
    test_loop(test_loader, model, loss)

print("Done!")

# Save the model weights, so it can be used for predictions.
torch.save(model.state_dict(), "/script_sight/ai/model_weights.pth")
