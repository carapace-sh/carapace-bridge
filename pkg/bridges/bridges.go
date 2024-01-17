package bridges

func filter(m map[string]bool, filter ...[]string) map[string]bool {
	for _, f := range filter {
		for _, e := range f {
			delete(m, e)
		}
	}
	return m
}
