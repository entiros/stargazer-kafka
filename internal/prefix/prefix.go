package prefix

import (
	"fmt"
	"strings"
)

// Validate prefix.
// Prefix example e85da0fd6.
func Validate(prefix string) error {

	if len(prefix) != 10 {
		return fmt.Errorf("prefix length must be 9. Found a prefix of length %d", len(prefix))
	}

	if !strings.HasPrefix(prefix, "e") {
		return fmt.Errorf("prefix must start with 'e'")
	}

	if !strings.HasSuffix(prefix, ".") {
		return fmt.Errorf("prefix must end with '.'")
	}

	return nil
}

func ExtractPrefix(topic string) (string, error) {

	err := Validate(topic[0:10])
	if err != nil {
		return "", err
	}
	return topic[0:10], nil

}

func HasPrefix(topic string) bool {

	err := Validate(topic[0:10])
	if err != nil {
		return true
	}

	return false

}
