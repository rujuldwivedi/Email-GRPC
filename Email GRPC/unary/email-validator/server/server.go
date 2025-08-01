package main

import (
    "context"
    "log"
    "net"
    "strings"

    "github.com/badoux/checkmail"
    "google.golang.org/grpc"
    "github.com/rujuldwivedi/Email-GRPC/proto/emailpb"
)

type server struct
{
    emailpb.UnimplementedEmailServiceServer
}

func (s *server) ValidateEmail(ctx context.Context, req *emailpb.EmailRequest) (*emailpb.EmailResponse, error)
{
    email := req.GetEmail()
    normalized := strings.ToLower(strings.TrimSpace(email))
    domainParts := strings.Split(normalized, "@")

    var suggested string
    if len(domainParts) == 2 {
        domain := domainParts[1]
        switch domain {
        case "gmal.com":
            suggested = "gmail.com"
        case "hotmial.com":
            suggested = "hotmail.com"
        }
    }

    err := checkmail.ValidateFormat(normalized)
    isValid := err == nil

    return &emailpb.EmailResponse{
        IsValid:         isValid,
        NormalizedEmail: normalized,
        SuggestedDomain: suggested,
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    emailpb.RegisterEmailServiceServer(grpcServer, &server{})

    log.Println("Server listening on :50051")
	
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
