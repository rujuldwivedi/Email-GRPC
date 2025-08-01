package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"github.com/rujuldwivedi/Email-GRPC/proto/emailpb"
)

type server struct {
	emailpb.UnimplementedEmailServiceServer
}

func (s *server) BatchValidate(stream emailpb.EmailService_BatchValidateServer) error {
	var responses []*emailpb.EmailResponse

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emailpb.BatchEmailResponse{Results: responses})
		}
		if err != nil {
			return err
		}

		email := req.GetEmail()
		resp := validateEmail(email)
		responses = append(responses, resp)
	}
}

func validateEmail(email string) *emailpb.EmailResponse {
	valid := strings.Contains(email, "@") && strings.Contains(email, ".")
	normal
