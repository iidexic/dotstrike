package cmd

func sliceUniques[E comparable](in []E) []E {
	seen := make(map[E]struct{}, len(in))
	out := make([]E, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func isEven(n int) bool { return n%2 == 0 }
