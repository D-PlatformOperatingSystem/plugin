// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "errors"

var (
	// ErrUnfreezeEmptied       
	ErrUnfreezeEmptied = errors.New("ErrUnfreezeEmptied")
	// ErrUnfreezeMeans        
	ErrUnfreezeMeans = errors.New("ErrUnfreezeMeans")
	// ErrUnfreezeID     ID  
	ErrUnfreezeID = errors.New("ErrUnfreezeID")
	// ErrNoPrivilege     
	ErrNoPrivilege = errors.New("ErrNoPrivilege")
	// ErrTerminated        
	ErrTerminated = errors.New("ErrTerminated")
)
