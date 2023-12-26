package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
)

type Cords struct {
	x float64
	y float64
}

// * Константы
const AccelerationOfGravity float64 = 9.8 // Ускорение свободного падения

const DiameterOfDisk float64 = 1.0 // Диаметр диска
const MassOfDisk float64 = 1.0     // Масса стержня

const DeltaT float64 = 0.001 // Шаг времени
const NumSteps int = 10000   // Количество шагов

const StepOfParya float64 = 0.01 // Шаг работника
const MassOfParya float64 = 0.1  // Масса работника

// * Слайс с позициями работников
var PositionsOfParyas = []Cords{}

// Братан, короче, эта функция делает уравнение
func Balance() {
	var results_X []float64
	var results_Y []float64

	for numOfWorkers := 1; numOfWorkers <= 24; numOfWorkers++ {
		// Инициализация начальных условий
		PositionsOfParyas = nil
		for i := 0; i < numOfWorkers; i++ {
			PositionsOfParyas = append(PositionsOfParyas, Cords{x: float64(rand.Intn(6))/10 - 0.25, y: float64(rand.Intn(6))/10 - 0.25})
		}

		theta_x := math.Pi / 9 // Начальный угол отклонения
		omega_x := 0.0         // Начальная угловая скорость

		theta_y := 0.0 // Начальный угол отклонения
		omega_y := 0.0 // Начальная угловая скорость

		count := 0
		var sliceOfTheta_x []float64
		var sliceOfTheta_y []float64

		// Цикл по времени
		for i := 0; i < NumSteps; i++ {

			// Вычисление момента силы
			var tau_x float64
			for j := 0; j < len(PositionsOfParyas); j++ {
				tau_x += MassOfParya*AccelerationOfGravity*DiameterOfDisk*math.Sin(theta_x) + MassOfParya*AccelerationOfGravity*PositionsOfParyas[j].x*math.Cos(theta_x)
			}

			var tau_y float64
			for j := 0; j < len(PositionsOfParyas); j++ {
				tau_y += MassOfParya*AccelerationOfGravity*DiameterOfDisk*math.Sin(theta_y) + MassOfParya*AccelerationOfGravity*PositionsOfParyas[j].y*math.Cos(theta_y)
			}

			// Обновление угла и угловой скорости
			// X
			alpha_x := tau_x / (MassOfDisk * DiameterOfDisk * DiameterOfDisk / 3.0) // Угловое ускорение
			omega_x += alpha_x * DeltaT

			theta_x += omega_x * DeltaT
			// Y
			alpha_y := tau_y / (MassOfDisk * DiameterOfDisk * DiameterOfDisk / 3.0) // Угловое ускорение
			omega_y += alpha_y * DeltaT

			theta_y += omega_y * DeltaT

			sliceOfTheta_x = append(sliceOfTheta_x, theta_x)
			sliceOfTheta_y = append(sliceOfTheta_y, theta_y)

			workerGroup := &sync.WaitGroup{}

			// Передвижение по оси х
			for index := range PositionsOfParyas {
				workerGroup.Add(1)

				go whereToGo(index, workerGroup, theta_x, theta_y, omega_x, omega_y)
			}

			workerGroup.Wait()

			// Остановка, если стержень достиг вертикального положения
			if math.Abs(theta_x) >= math.Pi/2.0 || math.Abs(theta_y) >= math.Pi/2.0 {
				break
			}
			count++
		}

		// var sum_X float64 = 0
		// for _, value_x := range sliceOfTheta_x {
		// 	sum_X += math.Abs(value_x)
		// }
		var sum_X float64 = 0
		for index := 0; index < NumSteps; index++ {
			if index < len(sliceOfTheta_x) {
				sum_X += math.Abs(sliceOfTheta_x[index])
			} else {
				sum_X += math.Pi / 2
			}

		}

		// var sum_Y float64 = 0
		// for _, value_y := range sliceOfTheta_y {
		// 	sum_Y += math.Abs(value_y)
		// }
		var sum_Y float64 = 0
		for index := 0; index < NumSteps; index++ {
			if index < len(sliceOfTheta_y) {
				sum_Y += math.Abs(sliceOfTheta_y[index])
			} else {
				sum_Y += math.Pi / 2
			}

		}

		results_X = append(results_X, sum_X*DeltaT)
		results_Y = append(results_Y, sum_Y*DeltaT)
		fmt.Println(numOfWorkers, count)
	}

	var ideal_X []float64
	var ideal_Y []float64
	var nums []int

	for index := range results_X {
		nums = append(nums, index+1)
		ideal_X = append(ideal_X, results_X[0]/float64(index+1))
		ideal_Y = append(ideal_Y, results_Y[0]/float64(index+1))
	}

	DrawChartWith2Lines(results_X, ideal_X, nums, "result1.html", "Рад * сек", "N", "Полученное значение по x", "Идевальное значение")
	DrawChartWith2Lines(results_Y, ideal_Y, nums, "result2.html", "Рад * сек", "N", "Полученное значение по y", "Идевальное значение")
	DrawChartWith2Lines(results_X, results_Y, nums, "result3.html", "Рад * сек", "N", "Полученное значение по x", "Полученное значение по y")
}

func whereToGo(index int, group *sync.WaitGroup, theta_x, theta_y, omega_x, omega_y float64) {
	defer group.Done()

	step_po_x := 0.0
	step_po_y := 0.0

	if theta_x > 0 && omega_x > -1 {
		if PositionsOfParyas[index].x-StepOfParya >= -DiameterOfDisk/2 {
			step_po_x -= StepOfParya
		}
	} else if theta_x < 0 && omega_x < 1 {
		if PositionsOfParyas[index].x+StepOfParya <= DiameterOfDisk/2 {
			step_po_x += StepOfParya
		}
	} else {
		if PositionsOfParyas[index].x < 0 {
			step_po_x += StepOfParya
		} else if PositionsOfParyas[index].x > 0 {
			step_po_x -= StepOfParya
		}
	}

	if theta_y > 0 && omega_y > -1 {
		if PositionsOfParyas[index].y-StepOfParya >= -DiameterOfDisk/2 {
			step_po_y -= StepOfParya
		}
	} else if theta_y < 0 && omega_y < 1 {
		if PositionsOfParyas[index].y+StepOfParya <= DiameterOfDisk/2 {
			step_po_y += StepOfParya
		}
	} else {
		if PositionsOfParyas[index].y < 0 {
			step_po_y += StepOfParya
		} else if PositionsOfParyas[index].y > 0 {
			step_po_y -= StepOfParya
		}
	}

	// if math.Pow(PositionsOfParyas[index].x+step_po_x, 2)+math.Pow(PositionsOfParyas[index].y+step_po_y, 2) > math.Pow(DiameterOfDisk/2, 2) {
	// 	step_po_x = 0
	// 	step_po_y = 0
	// }

	PositionsOfParyas[index].x += step_po_x
	PositionsOfParyas[index].y += step_po_y
}
