package utils

import (
	"math/big"
)

func FromRau(rau string, unit string) (string, error) {
	return convert(rau, unit, "div"), nil
}

func ToRau(num string, uint string) (string, error) {
	return convert(num, uint, "multipliedBy"), nil
}

func convert(num string, unit string, operator string) string  {
	switch unit {
	case "Rau":
		return num
	case "KRau":
		return big.NewInt(1000).String()
	case "MRau":
		return big.NewInt(1000000).String()
	case "GRau":
		return big.NewInt(1000000000).String()
	case "Qev":
		return big.NewInt(1000000000000).String()
	case "Jing":
		return big.NewInt(1000000000000000).String()
	default:
		return big.NewInt(1000000000000000000).String()
	}
}


