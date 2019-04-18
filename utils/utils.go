// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package utils

import "math/big"

const (
	// Rau is the smallest non-fungible token unit
	Rau int64 = 1
	// KRau is 1000 Rau
	KRau = Rau * 1000
	// MRau is 1000 KRau
	MRau = KRau * 1000
	// GRau is 1000 MRau
	GRau = MRau * 1000
	// Qev is 1000 GRau
	Qev = GRau * 1000
	// Jin is 1000 Qev
	Jin = Qev * 1000
	// Iotx is 1000 Jin, which should be fit into int64
	Iotx = Jin * 1000
)

func FromRau(rau, unit string) string {
	return convert(rau, unit, "div")
}
func ToRau(num, unit string) string {
	return convert(num, unit, "mul")
}
func convert(num, unit, operator string) string {
	numInt, ok := new(big.Int).SetString(num, 10)
	if !ok {
		return ""
	}
	unitInt := int64(1)
	switch unit {
	case "Rau":
		unitInt = Rau
	case "KRau":
		unitInt = KRau
	case "MRau":
		unitInt = MRau
	case "GRau":
		unitInt = GRau
	case "Qev":
		unitInt = Qev
	case "Jin":
		unitInt = Jin
	default:
		unitInt = Iotx
	}
	return bigOperator(numInt, unitInt, operator)
}
func bigOperator(numInt *big.Int, unit int64, operator string) string {
	if operator == "div" {
		return numInt.Div(numInt, big.NewInt(Iotx/unit)).Text(10)
	} else {
		return numInt.Mul(numInt, big.NewInt(unit)).Text(10)
	}
}
