package middleware

import (
	"context"
	"encoding/base64"
	"session-16-crud-user-docker-compose/config"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		// Check for public methods that do not require authentication
		publicMethods := []string{
			"/proto.user_service.v1.UserService/GetUsers",
			"/proto.user_service.v1.UserService/GetUserByID",
		}
		for _, method := range publicMethods {
			if info.FullMethod == method {
				return handler(ctx, req)
			}
		}

		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Metadata not provided")
		}

		// Get authorization header
		autHeader, ok := md["authorization"]
		if !ok || len(autHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization header is missing")
		}

		// Check for Basic Auth scheme
		if !strings.HasPrefix(autHeader[0], "Basic ") {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization scheme")
		}

		// Decode Base64 credentials
		decoded, err := base64.StdEncoding.DecodeString(autHeader[0][6:])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization token")
		}

		// Split the credentials into usernama and password
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization token")
		}

		username, password := creds[0], creds[1]

		// Validate the credentials
		if username != config.AuthBasicUsername || password != config.AuthBasicPassword {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid username and password")
		}

		return handler(ctx, req)
	}
}
