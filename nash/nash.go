package nash

import (
	"math"
	"parallel/dataSet"
	"parallel/efficiency"
)

// * Расчет критерия Нэша
func CalcNashCriterion(num int64, arrOfA []float64, arrOfB []dataSet.BWithIndexes) (nash float64) {
	nash = 0.0
	currentEfficiency := efficiency.CalcEfficiency(arrOfA, arrOfB, num)

	for i := 0; i < int(dataSet.NumOfAgents); i++ {
		if (num/int64(math.Pow(6, float64(i))))%6 != 5 {
			nash += math.Abs(efficiency.CalcEfficiency(arrOfA, arrOfB, num+int64(math.Pow(6, float64(i)))) - currentEfficiency)
		}
		if (num/int64(math.Pow(6, float64(i))))%6 != 0 {
			nash += math.Abs(efficiency.CalcEfficiency(arrOfA, arrOfB, num-int64(math.Pow(6, float64(i)))) - currentEfficiency)
		}
	}

	return nash
}
