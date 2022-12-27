package config

import (
	"testing"
)

func TestDiff(t *testing.T) {

	newList := []string{"X", "A", "B", "C"}

	old := []string{"A", "B", "E"}

	add, del := Diff(newList, old)

	t.Log(add)
	t.Log(del)

}
