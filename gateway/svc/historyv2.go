package svc

import (
	"context"
)

type HistoryV2Api interface {
	GetTransactionHistoryDetailV2(ctx context.Context, params interface{}) (interface{}, error)
	GetTransactionHistoryListV2(ctx context.Context, params interface{}) (interface{}, error)
}
