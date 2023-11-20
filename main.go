package main

import "fmt"

// Размер итогового квадратного изображения
const SIZE = 1000

const BLUE_IMAGE = "originals\\girl.jpg"
const RED_IMAGE = "originals\\cat.jpg"

func main() {
	fmt.Println("Start")

	WalkManyTimes()

	fmt.Println("Complete")

}
