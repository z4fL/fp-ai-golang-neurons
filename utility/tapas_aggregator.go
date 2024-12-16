package utility

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const ErrNoNumericValuesFound = "no numeric values found"

type TapasProcessor struct {
	Cells []string
}

func (tp *TapasProcessor) CountUniqueCells() (int, string) {
	result := make(map[string]bool)
	for _, room := range tp.Cells {
		result[room] = true
	}

	var uniqueValue []string
	for value := range result {
		uniqueValue = append(uniqueValue, value)
	}

	sort.Strings(uniqueValue) // Sort the unique values alphabetically

	var builder strings.Builder
	for i, value := range uniqueValue {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(value)
	}
	valueList := builder.String()

	return len(result), valueList // Number of unique items
}

func (tp *TapasProcessor) Sum() float64 {
	total := 0.0
	for _, cell := range tp.Cells {
		if num, err := strconv.ParseFloat(cell, 64); err == nil {
			total += num
		}
	}
	return total
}

func (tp *TapasProcessor) Average() float64 {
	total := 0.0
	count := 0
	for _, cell := range tp.Cells {
		if num, err := strconv.ParseFloat(cell, 64); err == nil {
			total += num
			count++
		} else {
			fmt.Printf("Error parsing cell value '%s': %v\n", cell, err)
		}
	}
	if count == 0 {
		return 0.0
	}
	return float64(total) / float64(count)
}

// Max returns the maximum numeric value from the Cells slice.
// If no numeric values are found, it returns an error.
func (tp *TapasProcessor) Max() (float64, error) {
	maxValue := math.Inf(-1)
	found := false
	for _, cell := range tp.Cells {
		if num, err := strconv.ParseFloat(cell, 64); err == nil {
			if !found || num > maxValue {
				maxValue = num
				found = true
			}
		}
	}
	if !found {
		return 0, fmt.Errorf(ErrNoNumericValuesFound)
	}
	return maxValue, nil
}

func (tp *TapasProcessor) Min() (float64, error) {
	minValue := math.Inf(1)
	found := false
	for _, cell := range tp.Cells {
		if num, err := strconv.ParseFloat(cell, 64); err == nil {
			if !found || num < minValue {
				minValue = num
				found = true
			}
		}
	}
	if !found {
		return 0, fmt.Errorf("no numeric values found")
	}
	return minValue, nil
}
