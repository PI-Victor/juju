package mstate_test

import (
	"bytes"
	. "launchpad.net/gocheck"
	"launchpad.net/mgo"
	"os/exec"
	"time"
)

type MgoSuite struct {
	Addr    string
	Session *mgo.Session
	output  bytes.Buffer
	server  *exec.Cmd
}

const mgoport = "27017" // 50017
const mgoaddr = "localhost:"+mgoport

func (s *MgoSuite) SetUpSuite(c *C) {
	mgo.SetStats(true)
	dbdir := c.MkDir()
	args := []string{
		"--dbpath", dbdir,
		"--bind_ip", "127.0.0.1",
		"--port", mgoport,
		"--nssize", "1",
		"--noprealloc",
		"--smallfiles",
		"--nojournal",
	}
	s.server = exec.Command("mongod", args...)
	err := s.server.Start()
	if err != nil {
		panic(err)
	}
}

func (s *MgoSuite) TearDownSuite(c *C) {
	s.server.Process.Kill()
	s.server.Process.Wait()
}

func (s *MgoSuite) SetUpTest(c *C) {
	err := DropAll(mgoaddr)
	if err != nil {
		panic(err)
	}
	mgo.ResetStats()
	s.Addr = mgoaddr
	s.Session, err = mgo.Dial(s.Addr)
	if err != nil {
		panic(err)
	}
}

func (s *MgoSuite) TearDownTest(c *C) {
	s.Session.Close()
	for i := 0; ; i++ {
		stats := mgo.GetStats()
		if stats.SocketsInUse == 0 && stats.SocketsAlive == 0 {
			break
		}
		if i == 20 {
			c.Fatal("Test left sockets in a dirty state")
		}
		c.Logf("Waiting for sockets to die: %d in use, %d alive", stats.SocketsInUse, stats.SocketsAlive)
		time.Sleep(5e8)
	}
}

func DropAll(mongourl string) (err error) {
	session, err := mgo.Dial(mongourl)
	if err != nil {
		return err
	}
	defer session.Close()

	names, err := session.DatabaseNames()
	if err != nil {
		return err
	}
	for _, name := range names {
		switch name {
		case "admin", "local", "config":
		default:
			err = session.DB(name).DropDatabase()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
