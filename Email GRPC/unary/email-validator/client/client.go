package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "github.com/rujuldwivedi/Email-GRPC/proto/emailpb"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Could not connect: %v", err)
    }
    defer conn.Close()

    client := emailpb.NewEmailServiceClient(conn)

    email := "  ExaMple@Gmal.com  "

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    res, err := client.ValidateEmail(ctx, &emailpb.EmailRequest{
        Email: email,
    })
    if err != nil {
        log.Fatalf("Error calling ValidateEmail: %v", err)
    }

    fmt.Printf("Original: %s\n", email)
    fmt.Printf("Normalized: %s\n", res.GetNormalizedEmail())
    fmt.Printf("Is Valid: %v\n", res.GetIsValid())
    fmt.Printf("Suggested Domain: %s\n", res.GetSuggestedDomain())
}
