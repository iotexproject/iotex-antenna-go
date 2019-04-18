// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package antenna

import (
	"github.com/iotexproject/iotex-antenna-go/iotx"
	"github.com/iotexproject/iotex-antenna-go/utils"
)

type Antenna struct {
	Iotx *iotx.Iotx
}

func NewAntenna(host string) (*Antenna, error) {
	iotx, err := iotx.NewIotx(host)
	if err != nil {
		return nil, err
	}
	antenna := &Antenna{}
	antenna.Iotx = iotx
	return antenna, nil
}
func FromRau(rau, unit string) string {
	return utils.FromRau(rau, unit)
}
func ToRau(num, unit string) string {
	return utils.ToRau(num, unit)
}
