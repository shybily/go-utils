package pool

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	usedAt    int64 // atomic
	conn      net.Conn
	createdAt time.Time

	rd *bufio.Reader
	wr *bufio.Writer

	pooled bool
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		conn:      conn,
		createdAt: time.Now(),
		rd:        bufio.NewReader(conn),
		wr:        bufio.NewWriter(conn),
	}
	return c
}

func (c *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&c.usedAt)
	return time.Unix(unix, 0)
}

func (c *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&c.usedAt, tm.Unix())
}

func (c *Conn) Created(t time.Duration) bool {
	return time.Now().Sub(c.createdAt) >= t
}

func (c *Conn) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Conn) Reader() *bufio.Reader {
	return c.rd
}

func (c *Conn) Writer() *bufio.Writer {
	return c.wr
}

func (c *Conn) WriteWithWriter(ctx context.Context, fn func(wr *bufio.Writer) error) error {
	if c.wr.Buffered() > 0 {
		return fmt.Errorf("conn bufferd")
	}
	if ctx != nil {
		if deadline, ok := ctx.Deadline(); ok {
			if err := c.SetWriteDeadline(deadline); err != nil {
				return err
			}
		}
	}
	if err := fn(c.wr); err != nil {
		return err
	}
	if c.wr.Buffered() > 0 {
		if err := c.wr.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) ReadWithReader(ctx context.Context, fn func(rd *bufio.Reader) error) error {
	if c.rd.Buffered() > 0 {
		return fmt.Errorf("conn bufferd")
	}
	if ctx != nil {
		if deadline, ok := ctx.Deadline(); ok {
			if err := c.SetReadDeadline(deadline); err != nil {
				return err
			}
		}
	}
	if err := fn(c.rd); err != nil {
		return err
	}
	if c.rd.Buffered() > 0 {
		c.rd.Reset(c.rd)
	}
	return nil
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Conn) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
