package greeter

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	apiv1 "github.com/damejeras/goose/api/gen/go/v1"
)

type Server struct{}

func (s *Server) SayHello(
	ctx context.Context,
	req *connect.Request[apiv1.SayHelloRequest],
) (*connect.Response[apiv1.SayHelloResponse], error) {
	res := connect.NewResponse(&apiv1.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	return res, nil
}
