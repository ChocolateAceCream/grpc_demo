package transport

import (
	"context"
	"sms/service" // Import the missing "service" package

	"github.com/go-kit/kit/endpoint"
)

type SMSEndpoint struct {
	SendSMS endpoint.Endpoint
}

func MakeSMSEndpoint(s service.SMSService) SMSEndpoint {
	return SMSEndpoint{
		SendSMS: makeSendSMSEndpoint(s),
	}
}

func makeSendSMSEndpoint(s service.SMSService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.SendSMSRequest)
		resp, err := s.SendSMS(ctx, req)
		return resp, err
	}
}
