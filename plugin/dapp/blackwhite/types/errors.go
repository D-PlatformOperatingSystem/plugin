// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "errors"

var (
	// ErrIncorrectStatus
	ErrIncorrectStatus = errors.New("ErrIncorrectStatus")
	// ErrRepeatPlayerAddr
	ErrRepeatPlayerAddr = errors.New("ErrRepeatPlayerAddress")
	// ErrNoTimeoutDone
	ErrNoTimeoutDone = errors.New("ErrNoTimeoutDone")
	// ErrNoExistAddr      ，
	ErrNoExistAddr = errors.New("ErrNoExistAddress")
	// ErrNoLoopSeq
	ErrNoLoopSeq = errors.New("ErrBlackwhiteFinalloopLessThanSeq")
)
