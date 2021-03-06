// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package raftlease_test

import (
	"bytes"
	"io"
	"time"

	"github.com/hashicorp/raft"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"

	"github.com/juju/juju/core/globalclock"
	"github.com/juju/juju/core/lease"
	"github.com/juju/juju/core/raftlease"
)

var zero time.Time

type fsmSuite struct {
	testing.IsolationSuite

	fsm *raftlease.FSM
}

var _ = gc.Suite(&fsmSuite{})

func (s *fsmSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)
	s.fsm = raftlease.NewFSM()
}

func (s *fsmSuite) apply(c *gc.C, command raftlease.Command) interface{} {
	data, err := command.Marshal()
	c.Assert(err, jc.ErrorIsNil)
	result := s.fsm.Apply(&raft.Log{Data: data})
	return result
}

func (s *fsmSuite) TestClaim(c *gc.C) {
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "me",
		Duration:  time.Second,
	}
	err := s.apply(c, command)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(s.fsm.Leases(zero), gc.DeepEquals,
		map[lease.Key]lease.Info{
			{"ns", "model", "lease"}: {
				Holder: "me",
				Expiry: offset(time.Second),
			},
		},
	)

	// Can't claim it again.
	err = s.apply(c, command)
	c.Assert(err, jc.Satisfies, lease.IsInvalid)

	// Someone else trying to claim the lease.
	command.Holder = "you"
	err = s.apply(c, command)
	c.Assert(err, jc.Satisfies, lease.IsInvalid)

}

func offset(d time.Duration) time.Time {
	return zero.Add(d)
}

func (s *fsmSuite) TestExtend(c *gc.C) {
	// Can't extend unless we've previously claimed.
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationExtend,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "me",
		Duration:  time.Second,
	}
	c.Assert(s.apply(c, command), jc.Satisfies, lease.IsInvalid)

	// Ok, so we'll claim it.
	command.Operation = raftlease.OperationClaim
	c.Assert(s.apply(c, command), jc.ErrorIsNil)

	// Now we can extend it.
	command.Operation = raftlease.OperationExtend
	command.Duration = 2 * time.Second
	c.Assert(s.apply(c, command), jc.ErrorIsNil)

	c.Assert(s.fsm.Leases(zero), gc.DeepEquals,
		map[lease.Key]lease.Info{
			{"ns", "model", "lease"}: {
				Holder: "me",
				Expiry: offset(2 * time.Second),
			},
		},
	)

	// Extending by a time less than the remaining duration doesn't
	// shorten the lease (but does succeed).
	command.Duration = time.Millisecond
	c.Assert(s.apply(c, command), jc.ErrorIsNil)

	c.Assert(s.fsm.Leases(zero), gc.DeepEquals,
		map[lease.Key]lease.Info{
			{"ns", "model", "lease"}: {
				Holder: "me",
				Expiry: offset(2 * time.Second),
			},
		},
	)

	// Someone else can't extend it.
	command.Holder = "you"
	c.Assert(s.apply(c, command), jc.Satisfies, lease.IsInvalid)
}

func (s *fsmSuite) TestExpire(c *gc.C) {
	// Can't expire a non-existent lease.
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationExpire,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
	}
	c.Assert(s.apply(c, command), jc.Satisfies, lease.IsInvalid)

	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "me",
		Duration:  time.Second,
	}), jc.ErrorIsNil)

	// Not allowed to expire too early.
	c.Assert(s.apply(c, command), jc.Satisfies, lease.IsInvalid)
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationSetTime,
		OldTime:   s.fsm.GlobalTime(),
		NewTime:   s.fsm.GlobalTime().Add(2 * time.Second),
	}), jc.ErrorIsNil)

	c.Assert(s.apply(c, command), jc.ErrorIsNil)
	c.Assert(s.fsm.Leases(zero), gc.DeepEquals, map[lease.Key]lease.Info{})
}

func (s *fsmSuite) TestSetTime(c *gc.C) {
	// Time always starts at 0.
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationSetTime,
		OldTime:   zero,
		NewTime:   zero.Add(2 * time.Second),
	}), jc.ErrorIsNil)
	c.Assert(s.fsm.GlobalTime(), gc.Equals, zero.Add(2*time.Second))

	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationSetTime,
		OldTime:   zero,
		NewTime:   zero.Add(time.Second),
	}), jc.Satisfies, globalclock.IsConcurrentUpdate)
}

func (s *fsmSuite) TestLeases(c *gc.C) {
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "me",
		Duration:  time.Second,
	}), jc.ErrorIsNil)
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns2",
		ModelUUID: "model2",
		Lease:     "lease",
		Holder:    "you",
		Duration:  4 * time.Second,
	}), jc.ErrorIsNil)

	c.Assert(s.fsm.Leases(zero), gc.DeepEquals,
		map[lease.Key]lease.Info{
			{"ns", "model", "lease"}: {
				Holder: "me",
				Expiry: offset(time.Second),
			},
			{"ns2", "model2", "lease"}: {
				Holder: "you",
				Expiry: offset(4 * time.Second),
			},
		},
	)
}

func (s *fsmSuite) TestApplyInvalidCommand(c *gc.C) {
	c.Assert(s.apply(c, raftlease.Command{
		Version:   300,
		Operation: raftlease.OperationSetTime,
		OldTime:   zero,
		NewTime:   zero.Add(2 * time.Second),
	}), jc.Satisfies, errors.IsNotValid)
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: "libera-me",
	}), jc.Satisfies, errors.IsNotValid)
}

func (s *fsmSuite) TestSnapshot(c *gc.C) {
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "me",
		Duration:  time.Second,
	}), jc.ErrorIsNil)
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationSetTime,
		OldTime:   zero,
		NewTime:   zero.Add(2 * time.Second),
	}), jc.ErrorIsNil)
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns2",
		ModelUUID: "model2",
		Lease:     "lease",
		Holder:    "you",
		Duration:  4 * time.Second,
	}), jc.ErrorIsNil)

	snapshot, err := s.fsm.Snapshot()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(snapshot, gc.DeepEquals, &raftlease.Snapshot{
		Version: 1,
		Entries: map[raftlease.SnapshotKey]raftlease.SnapshotEntry{
			{"ns", "model", "lease"}: {
				Holder:   "me",
				Start:    zero,
				Duration: time.Second,
			},
			{"ns2", "model2", "lease"}: {
				Holder:   "you",
				Start:    zero.Add(2 * time.Second),
				Duration: 4 * time.Second,
			},
		},
		GlobalTime: zero.Add(2 * time.Second),
	})
}

func (s *fsmSuite) TestRestore(c *gc.C) {
	c.Assert(s.apply(c, raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "ns",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "me",
		Duration:  time.Second,
	}), jc.ErrorIsNil)

	// Restoring overwrites the state.
	reader := closer{Reader: bytes.NewBuffer([]byte(snapshotYaml))}
	err := s.fsm.Restore(&reader)
	c.Assert(err, jc.ErrorIsNil)

	expected := &raftlease.Snapshot{
		Version: 1,
		Entries: map[raftlease.SnapshotKey]raftlease.SnapshotEntry{
			{"ns", "model", "lease"}: {
				Holder:   "me",
				Start:    zero,
				Duration: 5 * time.Second,
			},
			{"ns2", "model2", "lease"}: {
				Holder:   "you",
				Start:    zero.Add(2 * time.Second),
				Duration: 10 * time.Second,
			},
		},
		GlobalTime: zero.Add(3 * time.Second),
	}

	actual, err := s.fsm.Snapshot()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(actual, gc.DeepEquals, expected)
}

func (s *fsmSuite) TestSnapshotPersist(c *gc.C) {
	snapshot := &raftlease.Snapshot{
		Version: 1,
		Entries: map[raftlease.SnapshotKey]raftlease.SnapshotEntry{
			{"ns", "model", "lease"}: {
				Holder:   "me",
				Start:    zero,
				Duration: time.Second,
			},
			{"ns2", "model2", "lease"}: {
				Holder:   "you",
				Start:    zero.Add(2 * time.Second),
				Duration: 4 * time.Second,
			},
		},
		GlobalTime: zero.Add(2 * time.Second),
	}
	var buffer bytes.Buffer
	sink := fakeSnapshotSink{Writer: &buffer}
	err := snapshot.Persist(&sink)
	c.Assert(err, gc.ErrorMatches, "quam olim abrahe")
	c.Assert(sink.cancelled, gc.Equals, true)

	// Don't compare buffer bytes in output yaml directly, it's
	// dependent on map ordering.
	decoder := yaml.NewDecoder(&buffer)
	var loaded raftlease.Snapshot
	err = decoder.Decode(&loaded)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(&loaded, gc.DeepEquals, snapshot)
}

func (s *fsmSuite) TestCommandValidationExpire(c *gc.C) {
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationExpire,
		Namespace: "namespace",
		ModelUUID: "model",
		Lease:     "lease",
	}
	c.Assert(command.Validate(), gc.Equals, nil)
	command.Holder = "me"
	c.Assert(command.Validate(), gc.ErrorMatches, "expire with holder not valid")
	command.Holder = ""
	command.ModelUUID = ""
	c.Assert(command.Validate(), gc.ErrorMatches, "expire with empty model UUID not valid")
}

func (s *fsmSuite) TestCommandValidationClaim(c *gc.C) {
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationClaim,
		Namespace: "namespace",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "you",
		Duration:  time.Second,
	}
	c.Assert(command.Validate(), gc.Equals, nil)
	command.OldTime = time.Now()
	c.Assert(command.Validate(), gc.ErrorMatches, "claim with old time not valid")
	command.OldTime = time.Time{}
	command.Lease = ""
	c.Assert(command.Validate(), gc.ErrorMatches, "claim with empty lease not valid")
}

func (s *fsmSuite) TestCommandValidationExtend(c *gc.C) {
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationExtend,
		Namespace: "namespace",
		ModelUUID: "model",
		Lease:     "lease",
		Holder:    "you",
		Duration:  time.Second,
	}
	c.Assert(command.Validate(), gc.Equals, nil)
	command.NewTime = time.Now()
	c.Assert(command.Validate(), gc.ErrorMatches, "extend with new time not valid")
	command.OldTime = time.Time{}
	command.Namespace = ""
	c.Assert(command.Validate(), gc.ErrorMatches, "extend with empty namespace not valid")
}

func (s *fsmSuite) TestCommandValidationSetTime(c *gc.C) {
	command := raftlease.Command{
		Version:   1,
		Operation: raftlease.OperationSetTime,
		OldTime:   time.Now(),
		NewTime:   time.Now(),
	}
	c.Assert(command.Validate(), gc.Equals, nil)
	command.Duration = time.Minute
	c.Assert(command.Validate(), gc.ErrorMatches, "setTime with duration not valid")
	command.Duration = 0
	command.NewTime = time.Time{}
	c.Assert(command.Validate(), gc.ErrorMatches, "setTime with zero new time not valid")
}

type fakeSnapshotSink struct {
	io.Writer
	cancelled bool
}

func (s *fakeSnapshotSink) ID() string {
	return "fakeSink"
}

func (s *fakeSnapshotSink) Cancel() error {
	s.cancelled = true
	return nil
}

func (s *fakeSnapshotSink) Close() error {
	return errors.Errorf("quam olim abrahe")
}

type closer struct {
	io.Reader
	closed bool
}

func (c *closer) Close() error {
	c.closed = true
	return nil
}

var snapshotYaml = `
version: 1
entries:
  ? namespace: ns
    model-uuid: model
    lease: lease
  : holder: me
    start: 0001-01-01T00:00:00Z
    duration: 5s
  ? namespace: ns2
    model-uuid: model2
    lease: lease
  : holder: you
    start: 0001-01-01T00:00:02Z
    duration: 10s
global-time: 0001-01-01T00:00:03Z
`[1:]
