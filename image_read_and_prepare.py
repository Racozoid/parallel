import sys

from PIL import Image, ImageOps
from skimage.io import imread

SIZE: tuple[int, int] = (int(sys.argv[2]), int(sys.argv[2]))

with Image.open(sys.argv[1]) as image:
    ImageOps.grayscale(image).resize(SIZE).save("grayscale.jpg")

arr_of_pixels = imread("grayscale.jpg", True)

with open("list_of_pixels.txt", "w") as file:
    for y in arr_of_pixels:
        for x in y:
            file.write(str(x) + " ")
        file.write("\n")
