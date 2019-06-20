package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gsmcwhirter/go-util/v5/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v10/wsclient"
)

var errDone = errors.New("done")

type mockWSConn struct {
	doneFirst   bool
	idx         int
	deadlineSet bool
	deadline    time.Time

	msgType int
	first   [][]byte
	repeat  [][]byte
}

func (c *mockWSConn) Close() error { return nil }
func (c *mockWSConn) SetReadDeadline(t time.Time) error {
	c.deadlineSet = true
	c.deadline = t
	return nil
}
func (c *mockWSConn) ReadMessage() (int, []byte, error) {
	if c.deadlineSet && c.deadline.Before(time.Now()) {
		return 0, nil, errDone
	}

	if !c.doneFirst && c.idx >= len(c.first) {
		c.doneFirst = true
		c.idx = 0
	}

	if !c.doneFirst {
		m := c.first[c.idx]
		fmt.Println(m)
		c.idx++
		return c.msgType, m, nil
	}

	if c.idx >= len(c.repeat) {
		c.idx = 0
	}

	if len(c.repeat) == 0 {
		return 0, nil, errDone
	}

	m := c.repeat[c.idx]
	c.idx++
	fmt.Println(m)
	return c.msgType, m, nil
}
func (c *mockWSConn) WriteMessage(int, []byte) error {
	return nil
}

type mockWSDialer struct {
	MsgType int
	First   [][]byte
	Repeat  [][]byte
}

func (d *mockWSDialer) Dial(string, http.Header) (wsclient.Conn, *http.Response, error) {
	return &mockWSConn{
		msgType: d.MsgType,
		first:   d.First,
		repeat:  d.Repeat,
	}, &http.Response{StatusCode: 200}, nil
}
