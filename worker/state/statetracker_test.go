// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/state"
	statetesting "github.com/juju/juju/state/testing"
	workerstate "github.com/juju/juju/worker/state"
)

type StateTrackerSuite struct {
	statetesting.StateSuite
	pool         *state.State
	stateTracker workerstate.StateTracker
}

var _ = gc.Suite(&StateTrackerSuite{})

func (s *StateTrackerSuite) SetUpTest(c *gc.C) {
	s.StateSuite.SetUpTest(c)
	s.stateTracker = workerstate.NewStateTracker(s.StatePool)
}

func (s *StateTrackerSuite) TestDoneWithNoUse(c *gc.C) {
	err := s.stateTracker.Done()
	c.Assert(err, jc.ErrorIsNil)
	assertStatePoolClosed(c, s.StatePool)
}

func (s *StateTrackerSuite) TestTooManyDones(c *gc.C) {
	err := s.stateTracker.Done()
	c.Assert(err, jc.ErrorIsNil)
	assertStatePoolClosed(c, s.StatePool)

	err = s.stateTracker.Done()
	c.Assert(err, gc.Equals, workerstate.ErrStateClosed)
	assertStatePoolClosed(c, s.StatePool)
}

func (s *StateTrackerSuite) TestUse(c *gc.C) {
	pool, err := s.stateTracker.Use()
	c.Check(pool.SystemState(), gc.Equals, s.State)
	c.Check(err, jc.ErrorIsNil)

	pool, err = s.stateTracker.Use()
	c.Check(pool.SystemState(), gc.Equals, s.State)
	c.Check(err, jc.ErrorIsNil)
}

func (s *StateTrackerSuite) TestUseAndDone(c *gc.C) {
	// Ref count starts at 1 (the creator/owner)

	_, err := s.stateTracker.Use()
	// 2
	c.Check(err, jc.ErrorIsNil)

	_, err = s.stateTracker.Use()
	// 3
	c.Check(err, jc.ErrorIsNil)

	c.Check(s.stateTracker.Done(), jc.ErrorIsNil)
	// 2
	assertStatePoolNotClosed(c, s.StatePool)

	_, err = s.stateTracker.Use()
	// 3
	c.Check(err, jc.ErrorIsNil)

	c.Check(s.stateTracker.Done(), jc.ErrorIsNil)
	// 2
	assertStatePoolNotClosed(c, s.StatePool)

	c.Check(s.stateTracker.Done(), jc.ErrorIsNil)
	// 1
	assertStatePoolNotClosed(c, s.StatePool)

	c.Check(s.stateTracker.Done(), jc.ErrorIsNil)
	// 0
	assertStatePoolClosed(c, s.StatePool)
}

func (s *StateTrackerSuite) TestUseWhenClosed(c *gc.C) {
	c.Assert(s.stateTracker.Done(), jc.ErrorIsNil)

	pool, err := s.stateTracker.Use()
	c.Check(pool, gc.IsNil)
	c.Check(err, gc.Equals, workerstate.ErrStateClosed)
}

func assertStatePoolNotClosed(c *gc.C, pool *state.StatePool) {
	c.Assert(pool.SystemState(), gc.NotNil)
	err := pool.SystemState().Ping()
	c.Assert(err, jc.ErrorIsNil)
}

func assertStatePoolClosed(c *gc.C, pool *state.StatePool) {
	c.Assert(pool.SystemState(), gc.IsNil)
}
