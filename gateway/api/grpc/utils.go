package grpc

import (
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCConnRef struct {
	conn *grpc.ClientConn
	ref  int32
	mu   sync.Mutex
}

func NewGRPConnRef() *GRPCConnRef {
	return &GRPCConnRef{}
}

func (gcr *GRPCConnRef) Conn(instance string, logger log.Logger) (*grpc.ClientConn, error) {
	gcr.mu.Lock()
	defer gcr.mu.Unlock()
	level.Info(logger).Log("Connection GRPCConnRef With instance", instance)

	if gcr.ref == 0 {
		conn, err := grpc.NewClient(instance, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			level.Info(logger).Log("Error Connection GRPCConnRef", err)
			return nil, err
		}
		gcr.conn = conn
		level.Info(logger).Log("Status GRPCConnRef", err)
	}

	level.Info(logger).Log("Retry GRPCConnRef", gcr.ref)

	gcr.ref++
	return gcr.conn, nil
}

func (gcr *GRPCConnRef) Close() error {
	gcr.mu.Lock()
	defer gcr.mu.Unlock()

	gcr.ref--
	if gcr.ref <= 0 {
		return gcr.conn.Close()
	}
	return nil
}
