package executor

import "sort"

func Lamber(v uint32) uint32 {
	if v == 0 {
		return v
	}

	return v
}

func Sort(v []uint32) {
	go func() {
		u := Lamber(1) + uint32(1)

		_ = u
	}()
	sort.Slice(v, func(i, j int) bool {
		if v[i] == v[j] {
			return i < j
		}

		return v[i] < v[j]
	})
}
