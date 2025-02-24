package svc

//"context"

// Service interface
type Service interface {
	Auth() AuthApi
	Signature() SignatureApi
	Inquiry() InquiryApi
	History() HistoryApi
	Transfer() TransferApi
	Inquiryv2() InquiryV2Api
	Historyv2() HistoryV2Api
	Transferv2() TransferV2Api
}

var _ Service = (*service)(nil)

type service struct {
	authSvc       AuthApi
	signatureSvc  SignatureApi
	inquirySvc    InquiryApi
	historySvc    HistoryApi
	transferSvc   TransferApi
	inquiryv2Svc  InquiryV2Api
	historyv2Svc  HistoryV2Api
	transferv2Svc TransferV2Api
}

func (s *service) Auth() AuthApi {
	return s.authSvc
}

func (s *service) Signature() SignatureApi {
	return s.signatureSvc
}

func (s *service) Inquiry() InquiryApi {
	return s.inquirySvc
}

func (s *service) History() HistoryApi {
	return s.historySvc
}

func (s *service) Transfer() TransferApi {
	return s.transferSvc
}

func (s *service) Inquiryv2() InquiryV2Api {
	return s.inquiryv2Svc
}

func (s *service) Historyv2() HistoryV2Api {
	return s.historyv2Svc
}

func (s *service) Transferv2() TransferV2Api {
	return s.transferv2Svc
}

func NewService(authsvc AuthApi,
	signaturesvc SignatureApi,
	inquirysvc InquiryApi,
	historysvc HistoryApi,
	transfersvc TransferApi,
	inquiryv2svc InquiryApi,
	historyv2svc HistoryApi,
	transferv2svc TransferApi) Service {
	return &service{authsvc,
		signaturesvc,
		inquirysvc,
		historysvc,
		transfersvc,
		inquiryv2svc,
		historyv2svc,
		transferv2svc}
}
