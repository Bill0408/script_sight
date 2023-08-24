from torch import nn


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
