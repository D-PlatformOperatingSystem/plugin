// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//dpos action ty
const (
	DposVoteActionRegist = iota + 1
	DposVoteActionCancelRegist
	DposVoteActionReRegist
	DposVoteActionVote
	DposVoteActionCancelVote
	DposVoteActionRegistVrfM
	DposVoteActionRegistVrfRP
	DposVoteActionRecordCB
	DPosVoteActionRegistTopNCandidator

	CandidatorStatusRegist = iota + 1
	CandidatorStatusVoted
	CandidatorStatusCancelVoted
	CandidatorStatusCancelRegist
	CandidatorStatusReRegist

	VrfStatusMRegist = iota + 1
	VrfStatusRPRegist

	CBStatusRecord = iota + 1

	TopNCandidatorStatusRegist = iota + 1
)

//log ty
const (
	TyLogCandicatorRegist       = 1001
	TyLogCandicatorVoted        = 1002
	TyLogCandicatorCancelVoted  = 1003
	TyLogCandicatorCancelRegist = 1004
	TyLogCandicatorReRegist     = 1005
	TyLogVrfMRegist             = 1006
	TyLogVrfRPRegist            = 1007
	TyLogCBInfoRecord           = 1008
	TyLogTopNCandidatorRegist   = 1009
)

const (
	//VoteFrozenTime    = 3 * 24 * 3600

	//RegistFrozenCoins
	RegistFrozenCoins int64 = 1000000000000

	//VoteTypeNone
	VoteTypeNone int32 = 1

	//VoteTypeVote
	VoteTypeVote int32 = 2

	//VoteTypeCancelVote
	VoteTypeCancelVote int32 = 3

	//VoteTypeCancelAllVote
	VoteTypeCancelAllVote int32 = 4

	//TopNCandidatorsVoteInit topN    ：
	TopNCandidatorsVoteInit int64 = 0

	//TopNCandidatorsVoteMajorOK topN    ：2/3
	TopNCandidatorsVoteMajorOK int64 = 1

	//TopNCandidatorsVoteMajorFail topN    ：2/3
	TopNCandidatorsVoteMajorFail int64 = 2
)

//
//   github     ，        ,
//      ，
var (
	DPosX          = "dpos"
	ExecerDposVote = []byte(DPosX)
)

const (
	//FuncNameQueryCandidatorByPubkeys func name
	FuncNameQueryCandidatorByPubkeys = "QueryCandidatorByPubkeys"

	//FuncNameQueryCandidatorByTopN func name
	FuncNameQueryCandidatorByTopN = "QueryCandidatorByTopN"

	//FuncNameQueryVrfByTime func name
	FuncNameQueryVrfByTime = "QueryVrfByTime"

	//FuncNameQueryVrfByCycle func name
	FuncNameQueryVrfByCycle = "QueryVrfByCycle"

	//FuncNameQueryVrfByCycleForTopN func name
	FuncNameQueryVrfByCycleForTopN = "QueryVrfByCycleForTopN"

	//FuncNameQueryVrfByCycleForPubkeys func name
	FuncNameQueryVrfByCycleForPubkeys = "QueryVrfByCycleForPubkeys"

	//FuncNameQueryVote func name
	FuncNameQueryVote = "QueryVote"

	//CreateRegistTx
	CreateRegistTx = "Regist"

	//CreateCancelRegistTx
	CreateCancelRegistTx = "CancelRegist"

	//CreateReRegistTx
	CreateReRegistTx = "ReRegist"

	//CreateVoteTx
	CreateVoteTx = "Vote"

	//CreateCancelVoteTx
	CreateCancelVoteTx = "CancelVote"

	//CreateRegistVrfMTx     Vrf M
	CreateRegistVrfMTx = "RegistVrfM"

	//CreateRegistVrfRPTx     Vrf R/P
	CreateRegistVrfRPTx = "RegistVrfRP"

	//CreateRecordCBTx     CB
	CreateRecordCBTx = "RecordCB"

	//QueryVrfByTime   time  Vrf
	QueryVrfByTime = 1

	//QueryVrfByCycle   cycle  Vrf
	QueryVrfByCycle = 2

	//QueryVrfByCycleForTopN   cycle    topN      Vrf
	QueryVrfByCycleForTopN = 3

	//QueryVrfByCycleForPubkeys   cycle    pubkey        Vrf
	QueryVrfByCycleForPubkeys = 4

	//FuncNameQueryCBInfoByCycle func name
	FuncNameQueryCBInfoByCycle = "QueryCBInfoByCycle"

	//FuncNameQueryCBInfoByHeight func name
	FuncNameQueryCBInfoByHeight = "QueryCBInfoByHeight"

	//FuncNameQueryCBInfoByHash func name
	FuncNameQueryCBInfoByHash = "QueryCBInfoByHash"

	//FuncNameQueryLatestCBInfoByHeight func name
	//FuncNameQueryLatestCBInfoByHeight = "QueryLatestCBInfoByHeight"

	//QueryCBInfoByCycle   cycle  cycle boundary
	QueryCBInfoByCycle = 1

	//QueryCBInfoByHeight   stopHeight  cycle boundary
	QueryCBInfoByHeight = 2

	//QueryCBInfoByHash   stopHash  cycle boundary
	QueryCBInfoByHash = 3

	//QueryLatestCBInfoByHeight   stopHeight  cycle boundary
	//QueryLatestCBInfoByHeight = 4

	//FuncNameQueryTopNByVersion func name
	FuncNameQueryTopNByVersion = "QueryTopNByVersion"
)
