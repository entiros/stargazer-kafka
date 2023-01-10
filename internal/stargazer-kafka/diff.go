package stargazer_kafka

import "sort"

// We are not supposed to understand this.
// Read here:
//  http://www.mlsite.net/blog/?p=2250
//  http://www.mlsite.net/blog/?p=2271

type tuple struct {
	position int
	value    string
}

// Produce two lists required to make source equal to target.
// - inserts, what need to be inserted into source
// - deletes, what need to be deleted from source
func ListDiff(source []string, target []string) ([]string, []string) {

	sort.Strings(source)
	sort.Strings(target)

	ins, del := syncActions(source, target)

	var inserts []string
	for _, r := range ins {
		inserts = append(inserts, r.value)
	}

	var deletes []string
	for _, d := range del {
		deletes = append(deletes, source[d])
	}

	return inserts, deletes
}

func syncActions(a []string, target []string) ([]tuple, []int) {

	var inserts []tuple
	var deletes []int
	x := 0
	y := 0

	for (x < len(a)) || (y < len(target)) {
		if y >= len(target) {
			deletes = append(deletes, x)
			x += 1
		} else if x >= len(a) {
			inserts = append(inserts, tuple{y, target[y]})
			y += 1
		} else if a[x] < target[y] {
			deletes = append(deletes, x)
			x += 1
		} else if a[x] > target[y] {
			inserts = append(inserts, tuple{y, target[y]})
			y += 1
		} else {
			x += 1
			y += 1
		}
	}
	return inserts, deletes
}
