/*
 * Copyright D-Platform Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package types

import "errors"

var (
	// OracleX oracle name
	OracleX = "oracle"
)

// oracle action type
const (
	ActionEventPublish = iota + 1 //
	ActionResultPrePublish
	ActionResultPublish
	ActionEventAbort
	ActionResultAbort
)

// oracle status
const (
	NoEvent = iota
	EventPublished
	EventAborted
	ResultPrePublished
	ResultAborted
	ResultPublished
)

// log type define
const (
	TyLogEventPublish     = 810
	TyLogEventAbort       = 811
	TyLogResultPrePublish = 812
	TyLogResultAbort      = 813
	TyLogResultPublish    = 814
)

// executor action and function define
const (
	// FuncNameQueryOracleListByIDs   ids  OracleStatus
	FuncNameQueryOracleListByIDs = "QueryOraclesByIDs"
	// FuncNameQueryEventIDByStatus       eventID
	FuncNameQueryEventIDByStatus = "QueryEventIDsByStatus"
	// FuncNameQueryEventIDByAddrAndStatus             eventID
	FuncNameQueryEventIDByAddrAndStatus = "QueryEventIDsByAddrAndStatus"
	// FuncNameQueryEventIDByTypeAndStatus            eventID
	FuncNameQueryEventIDByTypeAndStatus = "QueryEventIDsByTypeAndStatus"
	// CreateEventPublishTx
	CreateEventPublishTx = "EventPublish"
	// CreateAbortEventPublishTx
	CreateAbortEventPublishTx = "EventAbort"
	// CreatePrePublishResultTx
	CreatePrePublishResultTx = "ResultPrePublish"
	// CreateAbortResultPrePublishTx
	CreateAbortResultPrePublishTx = "ResultAbort"
	// CreateResultPublishTx
	CreateResultPublishTx = "ResultPublish"
)

// query param define
const (
	// ListDESC
	ListDESC = int32(0)
	// DefaultCount
	DefaultCount = int32(20)
)

// Errors for oracle
var (
	ErrTimeMustBeFuture           = errors.New("ErrTimeMustBeFuture")
	ErrNoPrivilege                = errors.New("ErrNoPrivilege")
	ErrOracleRepeatHash           = errors.New("ErrOracleRepeatHash")
	ErrEventIDNotFound            = errors.New("ErrEventIDNotFound")
	ErrEventAbortNotAllowed       = errors.New("ErrEventAbortNotAllowed")
	ErrResultPrePublishNotAllowed = errors.New("ErrResultPrePublishNotAllowed")
	ErrPrePublishAbortNotAllowed  = errors.New("ErrPrePublishAbortNotAllowed")
	ErrResultPublishNotAllowed    = errors.New("ErrResultPublishNotAllowed")
	ErrParamNeedIDs               = errors.New("ErrParamNeedIDs")
	ErrParamStatusInvalid         = errors.New("ErrParamStatusInvalid")
	ErrParamAddressMustnotEmpty   = errors.New("ErrParamAddressMustnotEmpty")
	ErrParamTypeMustNotEmpty      = errors.New("ErrParamTypeMustNotEmpty")
)
