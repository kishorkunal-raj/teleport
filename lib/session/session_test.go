/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package session

import (
	"context"
	"testing"
	"time"

	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/backend/lite"
	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/utils"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	. "gopkg.in/check.v1"
)

func TestSessions(t *testing.T) { TestingT(t) }

type SessionSuite struct {
	dir   string
	srv   *server
	bk    backend.Backend
	clock clockwork.FakeClock
}

var _ = Suite(&SessionSuite{})

func (s *SessionSuite) SetUpSuite(c *C) {
	utils.InitLoggerForTests(testing.Verbose())
}

func (s *SessionSuite) SetUpTest(c *C) {
	var err error

	s.clock = clockwork.NewFakeClockAt(time.Date(2016, 9, 8, 7, 6, 5, 0, time.UTC))
	s.dir = c.MkDir()

	s.bk, err = lite.NewWithConfig(context.TODO(),
		lite.Config{
			Path:  s.dir,
			Clock: s.clock,
		},
	)
	c.Assert(err, IsNil)

	srv, err := New(s.bk)
	srv.(*server).clock = s.clock
	s.srv = srv.(*server)
	c.Assert(err, IsNil)
}

func (s *SessionSuite) TearDownTest(c *C) {
	c.Assert(s.bk.Close(), IsNil)
}

func (s *SessionSuite) TestID(c *C) {
	id := NewID()
	id2, err := ParseID(id.String())
	c.Assert(err, IsNil)
	c.Assert(id, Equals, *id2)

	for _, val := range []string{"garbage", "", "   ", string(id) + "extra"} {
		id := ID(val)
		c.Assert(id.Check(), NotNil)
	}
}

func (s *SessionSuite) TestSessionsCRUD(c *C) {
	out, err := s.srv.GetSessions(defaults.Namespace)
	c.Assert(err, IsNil)
	c.Assert(len(out), Equals, 0)

	// Create session.
	sess := Session{
		ID:             NewID(),
		Namespace:      defaults.Namespace,
		TerminalParams: TerminalParams{W: 100, H: 100},
		Login:          "bob",
		LastActive:     s.clock.Now().UTC(),
		Created:        s.clock.Now().UTC(),
	}
	c.Assert(s.srv.CreateSession(sess), IsNil)

	// Make sure only one session exists.
	out, err = s.srv.GetSessions(defaults.Namespace)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []Session{sess})

	// Make sure the session is the one created above.
	s2, err := s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(err, IsNil)
	c.Assert(s2, DeepEquals, &sess)

	// Update session terminal parameter
	err = s.srv.UpdateSession(UpdateRequest{
		ID:             sess.ID,
		Namespace:      defaults.Namespace,
		TerminalParams: &TerminalParams{W: 101, H: 101},
	})
	c.Assert(err, IsNil)

	// Verify update was applied.
	sess.TerminalParams = TerminalParams{W: 101, H: 101}
	s2, err = s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(err, IsNil)
	c.Assert(s2, DeepEquals, &sess)

	// Remove the session.
	err = s.srv.DeleteSession(defaults.Namespace, sess.ID)
	c.Assert(err, IsNil)

	// Make sure session no longer exists.
	_, err = s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(err, NotNil)
}

// TestSessionsInactivity makes sure that session will be marked
// as inactive after period of inactivity
func (s *SessionSuite) TestSessionsInactivity(c *C) {
	sess := Session{
		ID:             NewID(),
		Namespace:      defaults.Namespace,
		TerminalParams: TerminalParams{W: 100, H: 100},
		Login:          "bob",
		LastActive:     s.clock.Now().UTC(),
		Created:        s.clock.Now().UTC(),
	}
	c.Assert(s.srv.CreateSession(sess), IsNil)

	// move forward in time:
	s.clock.Advance(defaults.ActiveSessionTTL + time.Second)

	// should not be in active sessions:
	s2, err := s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(err, NotNil)
	c.Assert(trace.IsNotFound(err), Equals, true)
	c.Assert(s2, IsNil)
}

func (s *SessionSuite) TestPartiesCRUD(c *C) {
	// create session:
	sess := Session{
		ID:             NewID(),
		Namespace:      defaults.Namespace,
		TerminalParams: TerminalParams{W: 100, H: 100},
		Login:          "vincent",
		LastActive:     s.clock.Now().UTC(),
		Created:        s.clock.Now().UTC(),
	}
	c.Assert(s.srv.CreateSession(sess), IsNil)
	// add two people:
	parties := []Party{
		{
			ID:         NewID(),
			RemoteAddr: "1_remote_addr",
			User:       "first",
			ServerID:   "luna",
			LastActive: s.clock.Now().UTC(),
		},
		{
			ID:         NewID(),
			RemoteAddr: "2_remote_addr",
			User:       "second",
			ServerID:   "luna",
			LastActive: s.clock.Now().UTC(),
		},
	}
	err := s.srv.UpdateSession(UpdateRequest{
		ID:        sess.ID,
		Namespace: defaults.Namespace,
		Parties:   &parties,
	})
	c.Assert(err, IsNil)
	// verify they're in the session:
	copy, err := s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(err, IsNil)
	c.Assert(len(copy.Parties), Equals, 2)

	// empty update (list of parties must not change)
	err = s.srv.UpdateSession(UpdateRequest{ID: sess.ID, Namespace: defaults.Namespace})
	c.Assert(err, IsNil)
	copy, _ = s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(len(copy.Parties), Equals, 2)

	// remove the 2nd party:
	deleted := copy.RemoveParty(parties[1].ID)
	c.Assert(deleted, Equals, true)
	err = s.srv.UpdateSession(UpdateRequest{ID: copy.ID, Parties: &copy.Parties, Namespace: defaults.Namespace})
	c.Assert(err, IsNil)
	copy, _ = s.srv.GetSession(defaults.Namespace, sess.ID)
	c.Assert(len(copy.Parties), Equals, 1)

	// we still have the 1st party in:
	c.Assert(parties[0].ID, Equals, copy.Parties[0].ID)
}
