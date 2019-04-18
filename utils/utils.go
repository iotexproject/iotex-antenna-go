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
	switch unit {
	case "Rau":
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/Rau)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(Rau)).Text(10)
		}
	case "KRau":
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/KRau)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(KRau)).Text(10)
		}
	case "MRau":
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/MRau)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(MRau)).Text(10)
		}
	case "GRau":
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/GRau)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(GRau)).Text(10)
		}
	case "Qev":
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/Qev)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(Qev)).Text(10)
		}
	case "Jin":
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/Jin)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(Jin)).Text(10)
		}
	default:
		if operator == "div" {
			return numInt.Div(numInt, big.NewInt(Iotx/Iotx)).Text(10)
		} else {
			return numInt.Mul(numInt, big.NewInt(Iotx)).Text(10)
		}
	}
}

//// FromRau is a function to convert Rau to Iotx.
//func FromRau(rau *big.Int) int64 {
//	rau = rau.Div(rau, big.NewInt(1e18))
//	return rau.Int64()
//}
//
//// ToRau is a function to convert various units to Rau.
//func ToRau(iotx int64) *big.Int {
//	itx := big.NewInt(iotx)
//	return itx.Mul(itx, big.NewInt(1e18))
//}
