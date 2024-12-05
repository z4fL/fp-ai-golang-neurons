package utility

import (
	"fmt"
	"strconv"
	"strings"
)

type TapasProcessor struct {
	Cells *[]string
}

func (tp *TapasProcessor) Count() (int, string) {
	result := make(map[string]bool)
	for _, room := range *tp.Cells {
		result[room] = true
	}

	var uniqueValue []string
	for value := range result {
		uniqueValue = append(uniqueValue, value)
	}

	valueList := strings.Join(uniqueValue, ", ")

	return len(result), valueList // Jumlah unique items
}

func (tp *TapasProcessor) Sum() float64 {
	total := 0.0
	for _, cell := range *tp.Cells {
		if num, err := strconv.ParseFloat(cell, 64); err == nil {
			total += num
		}
	}
	return total
}

func (tp *TapasProcessor) Average() float64 {
	total := 0
	count := 0
	for _, cell := range *tp.Cells {
		if num, err := strconv.Atoi(cell); err == nil {
			total += num
			count++
		}
	}
	if count == 0 {
		return 0.0
	}
	return float64(total) / float64(count)
}

func (tp *TapasProcessor) Max() (int, error) {
	maxValue := 0
	found := false
	for _, cell := range *tp.Cells {
		if num, err := strconv.Atoi(cell); err == nil {
			if !found || num > maxValue {
				maxValue = num
				found = true
			}
		}
	}
	if !found {
		return 0, fmt.Errorf("no numeric values found")
	}
	return maxValue, nil
}

func (tp *TapasProcessor) Min() (int, error) {
	minValue := 0
	found := false
	for _, cell := range *tp.Cells {
		if num, err := strconv.Atoi(cell); err == nil {
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
