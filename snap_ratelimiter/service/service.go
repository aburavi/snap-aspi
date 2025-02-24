package service

import (
	"context"

	"snap_ratelimiter/models"

	"github.com/segmentio/ksuid"
)

type ObjectService interface {
	Method(ctx context.Context, req models.Trx) (interface{}, error)
}

type objectService struct {
}

func NewObjectService() ObjectService {
	return &objectService{}
}

func (o *objectService) Method(ctx context.Context, req models.Trx) (interface{}, error) {
	// log.Println("api hit => ")
	resp := models.StreamData{
		UUID:             ksuid.New().String(),
		TrxType:          req.TrxType,
		SystemTraceAudit: req.SystemTraceAudit,
		IData:            req,
	}

	return resp, nil

}
