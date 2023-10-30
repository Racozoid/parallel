package dataSet

//* Колличество агентов
const NumOfAgents int64 = 24

// Получить значение для
func GetAgentValue(value string) int64 {
	switch value {
	case "0":
		return -3
	case "1":
		return -2
	case "2":
		return -1
	case "3":
		return 1
	case "4":
		return 2
	case "5":
		return 3
	default:
		return 0
	}
}
