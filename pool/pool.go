package pool

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Options struct {
	Dialer  func(ctx context.Context) (net.Conn, error)
	OnClose func(*Conn) error

	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

type ConnPool struct {
	opt *Options

	dialErrorsNum uint32 // atomic

	lastDialError atomic.Value

	mu           sync.Mutex
	conns        []*Conn
	idleConns    []*Conn
	connsLen     int
	idleConnsLen int

	_closed uint32
	closeCh chan struct{}
}

func NewConnPool(opt *Options) *ConnPool {
	p := &ConnPool{
		opt:       opt,
		mu:        sync.Mutex{},
		conns:     make([]*Conn, 0, opt.PoolSize),
		idleConns: make([]*Conn, 0, opt.PoolSize),
		closeCh:   make(chan struct{}),
	}

	p.mu.Lock()
	p.initIdleConn()
	p.mu.Unlock()

	return p
}

func (p *ConnPool) initIdleConn() {
	if p.opt.MinIdleConns == 0 {
		return
	}
	for p.connsLen < p.opt.PoolSize && p.idleConnsLen < p.opt.MinIdleConns {
		p.connsLen++
		p.idleConnsLen++
		go func() {
			err := p.newIdleConn()
			if err != nil {
				p.mu.Lock()
				p.connsLen--
				p.idleConnsLen--
				p.mu.Unlock()
			}
		}()
	}
}

func (p *ConnPool) newIdleConn() error {
	cn, err := p.dialConn(context.TODO())
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	p.mu.Unlock()
	return nil
}

func (p *ConnPool) NewConn(ctx context.Context) (*Conn, error) {
	return p.newConn(ctx)
}

func (p *ConnPool) newConn(ctx context.Context) (*Conn, error) {
	conn, err := p.dialConn(ctx)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	p.conns = append(p.conns, conn)
	p.mu.Unlock()
	return conn, nil
}

func (p *ConnPool) dialConn(ctx context.Context) (*Conn, error) {
	if p.closed() {
		return nil, fmt.Errorf("closed")
	}

	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.opt.PoolSize) {
		return nil, p.getLastDialError()
	}

	netConn, err := p.opt.Dialer(ctx)
	if err != nil {
		p.setLastDialError(err)
		return nil, err
	}

	cn := NewConn(netConn)
	return cn, nil
}

func (p *ConnPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

func (p *ConnPool) setLastDialError(err error) {
	p.lastDialError.Store(err)
}

func (p *ConnPool) getLastDialError() error {
	err, _ := p.lastDialError.Load().(error)
	if err != nil {
		return err
	}
	return nil
}

func (p *ConnPool) Close() error {
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return fmt.Errorf("has benn closed")
	}
	close(p.closeCh)

	var firstErr error
	p.mu.Lock()
	for _, cn := range p.conns {
		if err := p.closeConn(cn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.connsLen = 0
	p.idleConns = nil
	p.idleConnsLen = 0
	p.mu.Unlock()

	return firstErr
}

func (p *ConnPool) closeConn(cn *Conn) error {
	if p.opt.OnClose != nil {
		_ = p.opt.OnClose(cn)
	}
	return cn.Close()
}
