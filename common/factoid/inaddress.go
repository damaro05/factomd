// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Input object for a factoid transaction.   contains an amount
// and the destination address.

package factoid

import (
	"bytes"
	. "github.com/FactomProject/factomd/common/interfaces"
	. "github.com/FactomProject/factomd/common/primitives"
)

type InAddress struct {
	TransAddress
}

var _ IInAddress = (*InAddress)(nil)

func (e *InAddress) JSONByte() ([]byte, error) {
	return EncodeJSON(e)
}

func (e *InAddress) JSONString() (string, error) {
	return EncodeJSONString(e)
}

func (e *InAddress) JSONBuffer(b *bytes.Buffer) error {
	return EncodeJSONToBuffer(e, b)
}

func (b InAddress) String() string {
	txt, err := b.CustomMarshalText()
	if err != nil {
		return "<error>"
	}
	return string(txt)
}

func (InAddress) GetDBHash() IHash {
	return Sha([]byte("InAddress"))
}

func (i InAddress) GetNewInstance() IBlock {
	return new(InAddress)
}

func (a InAddress) CustomMarshalText() (text []byte, err error) {
	return a.CustomMarshalText2("input")
}

/******************************
 * Helper functions
 ******************************/

func NewInAddress(address IAddress, amount uint64) IInAddress {
	oa := new(InAddress)
	oa.Amount = amount
	oa.Address = address
	return oa
}
