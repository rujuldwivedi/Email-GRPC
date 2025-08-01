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
	normalized := strings.ToLower(strings.TrimSpace(email))

	var suggestion string
	if strings.Contains(normalized, "gmal.com") {
		suggestion = "Did you mean gmail.com?"
	}

	return &emailpb.EmailResponse{
		IsValid:         valid,
		NormalizedEmail: normalized,
		SuggestedDomain: suggestion,
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	emailpb.RegisterEmailServiceServer(s, &server{})

	fmt.Println("🚀 Streaming Email Server running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
