package transport

import (
	"strconv"
	"strings"
)

// parseBudgetValue centralizes numeric HTTP budget parsing.
func parseBudgetValue(raw string) (int64, error) {
	return strconv.ParseInt(strings.TrimSuffix(raw, " "), 10, 64)
}
