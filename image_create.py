import sys
import numpy as np

from PIL import Image

image = []
with open(sys.argv[2], mode="r") as file:
    for line in file:
        image.append(line.split())

Image.fromarray(np.array(image, int)).convert("RGB").save(sys.argv[1])
