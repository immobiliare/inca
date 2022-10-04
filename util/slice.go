package util

func StringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func StringSliceDistinct(slice []string) (distinct []string) {
	var cache = make(map[string]bool)
	for _, entry := range slice {
		if _, ok := cache[entry]; !ok {
			cache[entry] = true
			distinct = append(distinct, entry)
		}
	}
	return
}

func StringSliceContains(slice []string, item string) bool {
	for _, sliceItem := range slice {
		if sliceItem == item {
			return true
		}
	}
	return false
}
