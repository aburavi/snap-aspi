package svc

import (
	"context"
)

type TransferApi interface {
	PostTransferInterBank(ctx context.Context, params interface{}) (interface{}, error)
	PostTransferIntraBank(ctx context.Context, params interface{}) (interface{}, error)
	PostPaymentHostToHost(ctx context.Context, params interface{}) (interface{}, error)
	PostTransferStatus(ctx context.Context, params interface{}) (interface{}, error)
}
