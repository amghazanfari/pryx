package utils

import "math/rand/v2"

func CreateRandomInt(startRange, endRange int) int {
	ranRange := endRange - startRange
	if ranRange <= 0 {
		panic("the range for random is non positive, clearly something is wrong")
	}
	return rand.IntN(ranRange)

}
