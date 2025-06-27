package helmutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFail(t *testing.T) {
	rule := Empty()
	rules := []string{`templates/.?*`, `tests/*.yaml`, `!tests/__snapshot__/`}
	for _, r := range rules {
		err := rule.parseRule(r)
		assert.NoError(t, err)
	}
	for _, p := range rule.patterns {
		fmt.Println(p)
	}
}
