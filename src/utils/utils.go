package utils

func CompareSlices(s1 []byte, s2 []byte) bool {
	if ((s1 == nil) || (s2 == nil)) || (len(s1) != len(s2)) {
		return false
	}

	for n := range s1 {
		if s1[n] != s2[n] {
			return false
		}
	}
	return true
}
