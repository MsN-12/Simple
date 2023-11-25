package gapi

import (
	"context"
	"fmt"
	"github.com/MsN-12/simpleBank/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationBearer = "bearer"
	authorizationHeader = "authorization"
)

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("unauthorized user")
	}

	return payload, nil
}
func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, accessibleRole := range accessibleRoles {
		if userRole == accessibleRole {
			return true
		}
	}
	return false
}
