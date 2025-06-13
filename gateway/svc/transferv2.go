package svc

import (
	"context"
)

type TransferV2Api interface {
	PostTransferInterBankV2(ctx context.Context, params interface{}) (interface{}, error)
	PostTransferIntraBankV2(ctx context.Context, params interface{}) (interface{}, error)
	PostPaymentHostToHostV2(ctx context.Context, params interface{}) (interface{}, error)
	PostTransferStatusV2(ctx context.Context, params interface{}) (interface{}, error)
}
