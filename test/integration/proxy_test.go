package integration

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/slcjordan/oops"
)

type TestTCPProxy struct {
	Addr string

	dest     string
	listener net.Listener
}

func NewTestTCPProxy(dest string, networkConditions []oops.Condition) *TestTCPProxy {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(fmt.Errorf("could not start tcp proxy: %w", err))
	}
	listener = oops.InjectListener(listener, networkConditions...)
	proxy := TestTCPProxy{
		Addr:     listener.Addr().String(),
		listener: listener,
		dest:     dest,
	}
	return &proxy
}

func (t *TestTCPProxy) Run(ctx context.Context) error {
	defer t.listener.Close()

	for {
		rw, err := t.listener.Accept()
		if err != nil {
			return err
		}
		go t.serve(ctx, rw)
	}
}

func (t *TestTCPProxy) serve(ctx context.Context, c net.Conn) {
	defer c.Close()

	d, err := net.Dial("tcp", t.dest)
	if err != nil {
		panic(fmt.Errorf("could not start tcp proxy: %w", err))
	}
	defer d.Close()

	//nolint:errcheck
	go io.Copy(c, d)
	//nolint:errcheck
	go io.Copy(d, c)
	<-ctx.Done()
}
