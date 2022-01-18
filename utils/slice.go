package utils

// AllButOne receives an array of strings as a parameter and returns the same
// array without one of its items, which is passed as the second parameter.
func AllButOne(items []string, item string) []string {
	result := []string{}
	for _, i := range items {
		if i == item {
			continue
		}
		result = append(result, i)
	}
	return result
}
