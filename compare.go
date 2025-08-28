package splunk

import "fmt"
import "go.k6.io/k6/js/modules"

func init() {
	modules.Register("k6/x/compare", new(Compare))
}

type Compare struct {
	ComparisonResult string
}

func (c *Compare) IsGreater(a, b int) bool {
	if a > b {
		c.ComparisonResult = fmt.Sprintf("%d is greater than %d", a, b)
		return true
	} else {
		c.ComparisonResult = fmt.Sprintf("%d is NOT greater than %d", a, b)
		return false
	}
}
