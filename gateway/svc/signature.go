package svc

import (
	"context"
)

type SignatureApi interface {
	PostAuthSignature(ctx context.Context, params interface{}) (interface{}, error)
	PostTrxSignature(ctx context.Context, params interface{}) (interface{}, error)
}
