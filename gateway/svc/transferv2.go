package svc

import (
	"context"
)

type TransferV2Api interface {
	PostTransferInterBank(ctx context.Context, params interface{}) (interface{}, error)
	PostTransferIntraBank(ctx context.Context, params interface{}) (interface{}, error)
	PostTransferStatus(ctx context.Context, params interface{}) (interface{}, error)
}
