package main

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestReadCsv(t *testing.T) {
	in := `1,37.966660,23.728308,1405594957
1,37.966627,23.728263,1405594966
`
	r := csv.NewReader(strings.NewReader(in))
	ReadCsv(r)
}
