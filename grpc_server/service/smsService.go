package service

import (
	"context"
	"net/http"

	"github.com/go-kit/log"
)

type SMSService struct {
	logger log.Logger
}

type SendSMSResponse struct {
	Code int32
	Resp string
}

type SendSMSRequest struct {
	Receiver string
	Message  string
}

func (s *SMSService) SendSMS(ctx context.Context, request SendSMSRequest) (r SendSMSResponse, err error) {
	// from := "+14696082336"
	// username := os.Getenv("TWILIO_USERNAME")
	// password := os.Getenv("TWILIO_PASSWORD")

	// client := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: username,
	// 	Password: password,
	// })

	// params := &twilioApi.CreateMessageParams{}
	// params.SetTo(request.Receiver)
	// params.SetFrom(from)
	// params.SetBody(request.Message)

	// resp, err := client.Api.CreateMessage(params)
	// if err != nil {
	// 	r.Code = http.StatusInternalServerError
	// 	r.Resp = err.Error()
	// } else {
	// 	response, _ := json.Marshal(*resp)
	// 	r.Code = http.StatusOK
	// 	r.Resp = string(response)
	// }

	// mock
	r.Code = http.StatusOK
	r.Resp = "success"

	return r, err
}

type Service interface {
	SendSMS(ctx context.Context, request SendSMSRequest) (resp SendSMSResponse, err error)
}

func NewSMSService(logger log.Logger) SMSService {
	return SMSService{
		logger: logger,
	}
}
