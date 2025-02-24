package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/aburavi/snaputils/proto/signature"

	"gateway/svc"
)

var _ svc.SignatureApi = (*ClientSignature)(nil)

type ClientSignature struct {
	postauthsignature endpoint.Endpoint
	posttrxsignature  endpoint.Endpoint
}

func (c *ClientSignature) PostAuthSignature(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.postauthsignature(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (c *ClientSignature) PostTrxSignature(ctx context.Context, params interface{}) (interface{}, error) {
	md := metadata.New(map[string]string{})
	data := params.(map[string]interface{})
	md.Append(data["header"].(string))
	md.Append(data["path"].(string))
	metadata.NewOutgoingContext(ctx, md)

	rsp, err := c.posttrxsignature(ctx, params)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewSignatureClient(instancer sd.Instancer, logger log.Logger) *ClientSignature {
	c := &ClientSignature{}
	var options []grpctransport.ClientOption
	options = append(options)

	c.postauthsignature = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Signature",
			"AuthSignature",
			encodeSignatureAuthRequest,
			decodeSignatureAuthResponse,
			signature.SignatureResponse{},
			options...,
		).Endpoint(), "Signature().PostAuthSignature"
	})

	c.posttrxsignature = GRPCClientEndpoint(logger, instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"Signature",
			"TrxSignature",
			encodeSignatureAuthRequest,
			decodeSignatureAuthResponse,
			signature.SignatureResponse{},
			options...,
		).Endpoint(), "Signature().PostTrxSignature"
	})

	return c
}
