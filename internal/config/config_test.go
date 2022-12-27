package config

import (
	"testing"
	"time"
)

func TestConfigWatcher(t *testing.T) {

	dir := "../../configs"

	configFiles, err := GetConfigs(dir)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := WatchConfigs(dir, func(newFiles []string) {
			t.Logf("Config files: %v", newFiles)

			addMe, deleteMe := Diff(newFiles, configFiles)

			for _, c := range deleteMe {
				DeleteConfig(c)
				deleteString(c, configFiles)
			}

			for _, c := range addMe {
				AddConfig(c)
				configFiles = append(configFiles, c)
			}

		})
		if err != nil {
			t.Log(err)
		}
	}()

	time.Sleep(20 * time.Second)

}

func TestDiff(t *testing.T) {

	newList := []string{"X", "A", "B", "C"}

	old := []string{"A", "B", "E"}

	add, del := Diff(newList, old)

	t.Log(add)
	t.Log(del)

}
