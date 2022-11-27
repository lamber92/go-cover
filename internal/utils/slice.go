package utils

import "strconv"

func StringsToInts(source []string) []int {
	result := make([]int, 0, len(source))
	for _, v := range source {
		i, err := strconv.Atoi(v)
		if err == nil {
			result = append(result, i)
		}
	}
	return result
}

func StringsToIntSet(source []string) map[int]struct{} {
	result := make(map[int]struct{}, len(source))
	for _, v := range source {
		i, err := strconv.Atoi(v)
		if err == nil {
			result[i] = struct{}{}
		}
	}
	return result
}
