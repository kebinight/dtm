/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package dtmcli

import (
	"github.com/dtm-labs/dtm/dtmcli/dtmimp"
)

// Saga struct of saga
type Saga struct {
	dtmimp.TransBase
	orders     map[int][]int
	concurrent bool
}

// NewSaga create a saga
func NewSaga(server string, gid string) *Saga {
	return &Saga{TransBase: *dtmimp.NewTransBase(gid, "saga", server, ""), orders: map[int][]int{}}
}

// Add add a saga step
func (s *Saga) Add(action string, compensate string, postData interface{}) *Saga {
	s.Steps = append(s.Steps, map[string]string{"action": action, "compensate": compensate})
	s.Payloads = append(s.Payloads, dtmimp.MustMarshalString(postData))
	return s
}

// AddBranchOrder specify that branch should be after preBranches. branch should is larger than all the element in preBranches
func (s *Saga) AddBranchOrder(branch int, preBranches []int) *Saga {
	s.orders[branch] = preBranches
	return s
}

// EnableConcurrent enable the concurrent exec of sub trans
func (s *Saga) EnableConcurrent() *Saga {
	s.concurrent = true
	return s
}

// Submit submit the saga trans
func (s *Saga) Submit() error {
	s.AddConcurrentContext()
	return dtmimp.TransCallDtm(&s.TransBase, s, "submit")
}

// AddConcurrentContext adds concurrent options to the request context
func (s *Saga) AddConcurrentContext() {
	if s.concurrent {
		s.CustomData = dtmimp.MustMarshalString(map[string]interface{}{"orders": s.orders, "concurrent": s.concurrent})
	}
}
