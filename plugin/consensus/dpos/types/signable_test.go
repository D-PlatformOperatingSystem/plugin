package types

import (
	"encoding/json"

	"github.com/D-PlatformOperatingSystem/dpos/common/crypto"
	"github.com/D-PlatformOperatingSystem/dpos/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
)

func init() {
	//    VRF，    SECP256K1
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic("init ConsensusCrypto failed.")
	}

	ConsensusCrypto = cr
}

func TestVote(t *testing.T) {
	filename := "./tmp_priv_validator.json"
	save(filename, privValidatorFile)
	privValidator := LoadOrGenPrivValidatorFS(filename)

	now := time.Now().Unix()
	//task := dpos.DecideTaskByTime(now)
	//  vote，   vote
	voteItem := &VoteItem{
		VotedNodeAddress: privValidator.Address,
		VotedNodeIndex:   int32(0),
		Cycle:            100,
		CycleStart:       18888,
		CycleStop:        28888,
		PeriodStart:      20000,
		PeriodStop:       21000,
		Height:           100,
	}
	encode, err := json.Marshal(voteItem)
	if err != nil {
		panic("Marshal vote failed.")
	}

	voteItem.VoteID = crypto.Ripemd160(encode)

	vote := &Vote{
		DPosVote: &DPosVote{
			VoteItem:         voteItem,
			VoteTimestamp:    now,
			VoterNodeAddress: privValidator.GetAddress(),
			VoterNodeIndex:   int32(0),
		},
	}
	assert.True(t, 0 == len(vote.Signature))

	chainID := "test-chain-Ep9EcD"
	privValidator.SignVote(chainID, vote)
	assert.True(t, 0 <= len(vote.Signature))
	vote2 := vote.Copy()
	err = vote2.Verify(chainID, privValidator.PubKey)
	require.Nil(t, err)
	assert.True(t, 0 < len(vote.Hash()))
	remove(filename)
}

func TestNotify(t *testing.T) {
	filename := "./tmp_priv_validator.json"
	save(filename, privValidatorFile)
	privValidator := LoadOrGenPrivValidatorFS(filename)

	now := time.Now().Unix()
	//task := dpos.DecideTaskByTime(now)
	//  vote，   vote
	voteItem := &VoteItem{
		VotedNodeAddress: privValidator.Address,
		VotedNodeIndex:   int32(0),
		Cycle:            100,
		CycleStart:       18888,
		CycleStop:        28888,
		PeriodStart:      20000,
		PeriodStop:       21000,
		Height:           100,
	}
	encode, err := json.Marshal(voteItem)
	if err != nil {
		panic("Marshal vote failed.")
	}

	voteItem.VoteID = crypto.Ripemd160(encode)

	chainID := "test-chain-Ep9EcD"

	notify := &Notify{
		DPosNotify: &DPosNotify{
			Vote:              voteItem,
			HeightStop:        200,
			HashStop:          []byte("abcdef121212"),
			NotifyTimestamp:   now,
			NotifyNodeAddress: privValidator.GetAddress(),
			NotifyNodeIndex:   int32(0),
		},
	}

	err = privValidator.SignNotify(chainID, notify)
	require.Nil(t, err)

	notify2 := notify.Copy()
	err = notify2.Verify(chainID, privValidator.PubKey)
	require.Nil(t, err)
	assert.True(t, 0 < len(notify.Hash()))
	remove(filename)
}
