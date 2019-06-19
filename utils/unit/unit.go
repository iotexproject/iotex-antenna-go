// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package unit

import "math/big"

// IotexUnit defines iotex unit type.
type IotexUnit int64

const (
	// Rau is the smallest non-fungible token unit
	Rau IotexUnit = 1
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

// FromString converts string to IotexUnit.
func (u *IotexUnit) FromString(s string) {
	var unitInt IotexUnit
	switch s {
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
	if u == nil {
		u = &unitInt
	} else {
		*u = unitInt
	}
}

// FromRau converts Rau string into diffrent unit string
func FromRau(rau *big.Int, unit string) *big.Int {
	n := big.NewInt(0).Set(rau)
	u := IotexUnit(0)
	u.FromString(unit)
	return n.Div(n, big.NewInt(int64(u)))
}

// ToRau converts different unit string into Rau string
func ToRau(num *big.Int, unit string) *big.Int {
	n := big.NewInt(0).Set(num)
	u := IotexUnit(0)
	u.FromString(unit)
	return n.Mul(n, big.NewInt(int64(u)))
}
