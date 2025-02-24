package grpc

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"

	"github.com/aburavi/snaputils/utils"
)

var GRPCConnections sync.Map

func newGRPCClientFactory(logger log.Logger, makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string)) sd.Factory {
	return func(instance string) (i endpoint.Endpoint, closer io.Closer, e error) {
		ref := NewGRPConnRef()
		actual, _ := GRPCConnections.LoadOrStore(instance, ref)
		ref = actual.(*GRPCConnRef)

		conn, err := ref.Conn(instance, logger)
		if err != nil {
			level.Info(logger).Log(err.Error())
			return nil, nil, err
		}
		ep, _ := makeEndpoint(conn)
		return ep, ref, nil
	}
}

func GRPCClientEndpoint(logger log.Logger, instancer sd.Instancer, makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string)) endpoint.Endpoint {
	factory := newGRPCClientFactory(logger, makeEndpoint)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	timeout := 10 * time.Second
	retryMax := 3

	return lb.RetryWithCallback(timeout, balancer, func(n int, received error) (bool, error) {
		if _, ok := received.(*utils.Error); ok {
			level.Info(logger).Log("connection received", strconv.FormatBool(ok))
			return false, received
		}

		level.Info(logger).Log("connection retry", received.Error())
		return n < retryMax, received
	})
}
