package efficiency

import (
	"parallel/dataSet"
	"strings"
)

func IntToBase6(n int64) string {
	const base = 6
	digits := "012345"

	if n == 0 {
		return strings.Repeat("0", int(dataSet.NumOfAgents))
	}

	s := ""
	for ; n > 0; n = n / base {
		s = string(digits[n%base]) + s
	}

	return strings.Repeat("0", int(dataSet.NumOfAgents)-len(s)) + s
}

func CalcEfficiency(arrOfA []float64, arrOfB []dataSet.BWithIndexes, num int64) (efficiency float64) {
	efficiency = 0.0
	numInSix := IntToBase6(num)

	for i, a := range arrOfA {
		efficiency += a * float64(dataSet.GetAgentValue(string(numInSix[i])))
	}

	for _, b := range arrOfB {
		efficiency += b.Value * float64(dataSet.GetAgentValue(string(numInSix[b.I]))) * float64(dataSet.GetAgentValue(string(numInSix[b.J])))
	}

	return efficiency
}
