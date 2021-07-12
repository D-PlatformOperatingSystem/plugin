// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "errors"

var (
	// ErrVotePeriod      
	ErrVotePeriod = errors.New("ErrVotePeriod")
	// ErrProposalStatus     
	ErrProposalStatus = errors.New("ErrProposalStatus")
	// ErrRepeatVoteAddr       
	ErrRepeatVoteAddr = errors.New("ErrRepeatVoteAddr")
	// ErrRevokeProposalPeriod        
	ErrRevokeProposalPeriod = errors.New("ErrRevokeProposalPeriod")
	// ErrRevokeProposalPower     
	ErrRevokeProposalPower = errors.New("ErrRevokeProposalPower")
	// ErrTerminatePeriod     
	ErrTerminatePeriod = errors.New("ErrTerminatePeriod")
	// ErrNoActiveBoard        
	ErrNoActiveBoard = errors.New("ErrNoActiveBoard")
	// ErrNoAutonomyExec  Autonomy   
	ErrNoAutonomyExec = errors.New("ErrNoAutonomyExec")
	// ErrNoPeriodAmount         
	ErrNoPeriodAmount = errors.New("ErrNoPeriodAmount")
	// ErrMinerAddr       
	ErrMinerAddr = errors.New("ErrMinerAddr")
	// ErrBindAddr       
	ErrBindAddr = errors.New("ErrBindAddr")
	// ErrChangeBoardAddr            
	ErrChangeBoardAddr = errors.New("ErrChangeBoardAddr")
	// ErrBoardNumber         
	ErrBoardNumber = errors.New("ErrBoardNumber")
	// ErrRepeatAddr     
	ErrRepeatAddr = errors.New("ErrRepeatAddr")
	// ErrNotEnoughFund     
	ErrNotEnoughFund = errors.New("ErrNotEnoughFund")
	// ErrSetBlockHeight block height not match
	ErrSetBlockHeight = errors.New("ErrSetBlockHeight")
)
