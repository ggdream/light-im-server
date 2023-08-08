package exist

func ExistInSlice[T comparable](target T, colltections []T) bool {
	for _, v := range colltections {
		if target == v {
			return true
		}
	}

	return false
}
