package cmd

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type lildeq[T any] struct {
	d []T
}

func getdeq[T any](s []T) lildeq[T] { return lildeq[T]{s} }

func (d *lildeq[T]) popbackN(n int) []T {
	ct := len(d.d)
	if n <= 0 || ct <= 0 {
		return []T{}
	}
	r := min(n, ct)
	out := make([]T, r)
	for i := range r {
		out[i] = d.popback()
	}
	return out
}
func (d *lildeq[T]) popfrontN(n int) []T {
	ct := len(d.d)
	if n <= 0 || ct <= 0 {
		return []T{}
	}
	r := min(n, ct)
	out := make([]T, r)
	for i := range r {
		out[i] = d.popfront()
	}
	return out
}
func (d *lildeq[T]) popback() T {
	defer func() {
		d.d = d.d[:len(d.d)-1]
	}()
	return d.d[len(d.d)-1]
}
func (d *lildeq[T]) popfront() T {
	defer func() {
		d.d = d.d[1:]
	}()
	return d.d[0]
}
func organizeArglist(args []string, segmentPattern []int) [][]string {
	dq := getdeq(args)
	arglists := make([][]string, len(segmentPattern))
	argN := len(args)
	for i, seg := range segmentPattern {
		switch {
		case argN == 0:
			arglists[i] = []string{}
		case abs(seg) >= argN:
			arglists[i] = args
		case seg < 0:
			arglists[i] = dq.popbackN(-seg)
			argN += seg
		case seg == 0:
			copy(arglists[i], dq.d)
			argN = 0
		case seg > 0:
			arglists[i] = dq.popfrontN(seg)
			argN -= seg

		}
	}
	return arglists
}
