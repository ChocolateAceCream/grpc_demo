package transport

import (
	"context"
	"errors"
	"sms/proto"
	"sms/service"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type gRPCServer struct {
	SendSMSHandler gt.Handler
	proto.UnimplementedSmsServiceServer
}

func (s *gRPCServer) SendSMS(ctx context.Context, req *proto.SendSMSRequest) (*proto.SendSMSResponse, error) {
	_, resp, err := s.SendSMSHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.SendSMSResponse), nil
}

func decodeSendSMSRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*proto.SendSMSRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}
	r := service.SendSMSRequest{
		Receiver: req.Receiver,
		Message:  req.Message,
	}
	return r, nil
}

func encodeSendSMSResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(service.SendSMSResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}
	r := &proto.SendSMSResponse{
		Code: resp.Code,
		Resp: resp.Resp,
	}
	return r, nil
}

func NewGRPCServer(endpoint SMSEndpoint, logger log.Logger) proto.SmsServiceServer {
	return &gRPCServer{
		SendSMSHandler: gt.NewServer(
			endpoint.SendSMS,
			decodeSendSMSRequest,
			encodeSendSMSResponse,
		),
	}
}
