package testHelper_test

import (
	. "github.com/FactomProject/factomd/testHelper"
	"testing"
)

/*
func TestTest(t *testing.T) {
	privKey, pubKey, add := NewFactoidAddressStrings(1)
	t.Errorf("%v, %v, %v", privKey, pubKey, add)
}
*/

func Test(t *testing.T) {
	ecBlock := CreateTestEntryCreditBlock(nil)
	t.Errorf("%v", ecBlock.String())
}
