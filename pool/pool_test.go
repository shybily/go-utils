package pool

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"
)

func TestNewConnPool(t *testing.T) {

	p := NewConnPool(&Options{
		Dialer: func(ctx context.Context) (net.Conn, error) {
			return net.Dial("tcp", "127.0.0.1:6379")
		},
		OnClose:            nil,
		PoolSize:           50,
		MinIdleConns:       10,
		MaxConnAge:         time.Minute * 30,
		PoolTimeout:        time.Minute,
		IdleTimeout:        time.Minute * 10,
		IdleCheckFrequency: time.Minute,
	})

	conn, err := p.NewConn(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	err = conn.WriteWithWriter(context.Background(), func(wr *bufio.Writer) error {
		wr.WriteString("info \n")
		return wr.Flush()
	})
	if err != nil {
		t.Fatal(err)
	}

	err = conn.ReadWithReader(context.TODO(), func(rd *bufio.Reader) error {
		resp, err := rd.ReadString('\n')
		if err != nil {
			return err
		}
		println(resp)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
	_ = p.Close()
}
