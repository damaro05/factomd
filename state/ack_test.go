// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package state_test

import (
	"testing"
	//"time"

	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/state"
	. "github.com/FactomProject/factomd/testHelper"
	//. "github.com/FactomProject/factomd/state"
)

func TestIsStateFullySynced(t *testing.T) {
	s1_good := CreateAndPopulateTestState()
	t.Log("IsStateFullySynced():", s1_good.IsStateFullySynced())
	if !s1_good.IsStateFullySynced() {
		t.Error("Test state is shown not to be fully synced")
	}

	//// we can't test the negative here because when we set the bad DBHeight the
	//// state.ValidatorLoop() will panic before we call IsStateFullySynced()
	//	s2_bad := CreateAndPopulateTestState()
	//// the next line will cause the ValidatorLoop to panic
	//	s2_bad.ProcessLists.DBHeightBase = s2_bad.ProcessLists.LastList().DBHeight+10
	//	t.Log("IsStateFullySynced:", s2_bad.IsStateFullySynced())

}

func TestFetchECTransactionByHash(t *testing.T) {
	s1 := CreateAndPopulateTestState()
	blocks := CreateFullTestBlockSet()

	for _, block := range blocks {
		for _, tx := range block.ECBlock.GetEntries() {
			if tx.ECID() != entryCreditBlock.ECIDChainCommit &&
				tx.ECID() != entryCreditBlock.ECIDEntryCommit ||
				tx.ECID() == entryCreditBlock.ECIDBalanceIncrease {
				continue
			}

			// get the tx from the database
			dtx, err := s1.FetchECTransactionByHash(tx.Hash())
			if err != nil {
				t.Error("Could not fetch transaction:", err)
			}
			if dtx == nil {
				t.Error("Transaction not found in database")
				continue
			}

			// test that the db transaction matches the tx we are looking for
			p1, err := tx.MarshalBinary()
			if err != nil {
				t.Error(err)
			}
			p2, err := dtx.MarshalBinary()
			if err != nil {
				t.Error(err)
			}
			if !primitives.AreBytesEqual(p1, p2) {
				t.Error("Database transaction does not match transaction")
			}
		}
	}
}

func TestFetchFactoidTransactionByHash(t *testing.T) {
	s1 := CreateAndPopulateTestState()
	blocks := CreateFullTestBlockSet()

	for _, block := range blocks {
		for _, tx := range block.FBlock.GetTransactions() {
			// get the transaction from the database
			dtx, err := s1.FetchFactoidTransactionByHash(tx.GetHash())
			if err != nil {
				t.Error("Could not fetch transaction:", err)
			}
			if dtx == nil {
				t.Error("Transaction not found in database")
				continue
			}

			// make sure the tx matches the one we are looking for
			p1, err := tx.MarshalBinary()
			if err != nil {
				t.Error(err)
			}
			p2, err := tx.MarshalBinary()
			if err != nil {
				t.Error(err)
			}
			if !primitives.AreBytesEqual(p1, p2) {
				t.Error("Database transaction does not match transaction")
			}
		}
	}
}

func TestFetchPaidFor(t *testing.T) {
	s1 := CreateAndPopulateTestState()
	blocks := CreateFullTestBlockSet()

	for _, block := range blocks {
		for _, tx := range block.ECBlock.GetEntries() {
			switch tx.ECID() {
			case entryCreditBlock.ECIDEntryCommit:
				// check that we can get the hash for the paid entry commit
				eh := tx.(*entryCreditBlock.CommitEntry).EntryHash
				h1, err := s1.FetchPaidFor(eh)
				if err != nil {
					t.Error("Transaction not found in database:", err)
					continue
				}
				if h1 == nil {
					t.Error("Transaction not found in database")
					continue
				}

				// make sure the tx sig matches the one we got
				if !h1.IsSameAs(tx.GetSigHash()) {
					t.Error("Hash mismatch")
				}
			case entryCreditBlock.ECIDChainCommit:
				// check that we can get the hash for the paid chain commit
				eh := tx.(*entryCreditBlock.CommitChain).EntryHash
				h1, err := s1.FetchPaidFor(eh)
				if err != nil {
					t.Error("Transaction not found in database:", err)
					continue
				}
				if h1 == nil {
					t.Error("Transaction not found in database")
					continue
				}

				// make sure the tx sig matches the one we got
				if !h1.IsSameAs(tx.GetSigHash()) {
					t.Error("Hash mismatch")
				}
			default:
				// make sure we dont get a positive result for a non-paid entry
				h1, err := s1.FetchPaidFor(tx.Hash())
				if err != nil {
					t.Error(err)
				}
				if h1 != nil {
					t.Error("Found non-paid transaction")
				}
			}
		}
	}
}

func TestFetchEntryByHash(t *testing.T) {
	s1 := CreateAndPopulateTestState()
	blocks := CreateFullTestBlockSet()

	for _, block := range blocks {
		for _, h := range block.EBlock.GetEntryHashes() {
			// get the entry from the database
			dentry, err := s1.FetchEntryByHash(h)
			if err != nil {
				t.Error("Could not fetch entry:", err)
			}
			if dentry == nil {
				t.Error("Entry not found in database")
				continue
			}

			// make sure the entry hash matches the one we are looking for
			if !h.IsSameAs(dentry.GetHash()) {
				t.Error("Mismatched entry")
			}
		}
	}
}

func newEOM(s *state.State) *messages.EOM {
	eom := new(messages.EOM)
	eom.Timestamp = primitives.NewTimestampFromMilliseconds(0xFF22100122FF)
	eom.Minute = 3
	eom.ChainID = s.IdentityChainID
	eom.DBHeight = s.LLeaderHeight

	return eom
}

func TestNewAck(t *testing.T) {
	s := CreateAndPopulateTestState()
	eom := newEOM(s)
	ackMsg := s.NewAck(eom, primitives.NewZeroHash())
	ack, ok := ackMsg.(*messages.Ack)
	if !ok {
		t.Error("NewAck() created a non-ack message")
	}
	if eom.DBHeight != ack.DBHeight {
		t.Errorf("EOM DBheight does not match ack DBHeight: %d %d\n", eom.DBHeight, ack.DBHeight)
	}
}
