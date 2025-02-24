package grpc

import (
	"gateway/svc"
)

var _ svc.Service = (*service)(nil)

type service struct {
	authSvc       svc.AuthApi
	signatureSvc  svc.SignatureApi
	inquirySvc    svc.InquiryApi
	historySvc    svc.HistoryApi
	transferSvc   svc.TransferApi
	inquiryv2Svc  svc.InquiryV2Api
	historyv2Svc  svc.HistoryV2Api
	transferv2Svc svc.TransferV2Api
}

func (s *service) Auth() svc.AuthApi {
	return s.authSvc
}

func (s *service) Signature() svc.SignatureApi {
	return s.signatureSvc
}

func (s *service) Inquiry() svc.InquiryApi {
	return s.inquirySvc
}

func (s *service) History() svc.HistoryApi {
	return s.historySvc
}

func (s *service) Transfer() svc.TransferApi {
	return s.transferSvc
}

func (s *service) Inquiryv2() svc.InquiryV2Api {
	return s.inquiryv2Svc
}

func (s *service) Historyv2() svc.HistoryV2Api {
	return s.historyv2Svc
}

func (s *service) Transferv2() svc.TransferV2Api {
	return s.transferv2Svc
}
