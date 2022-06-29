package investapi

import (
	"golang.org/x/oauth2"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func CreateInstrumentsServiceClient(url string, token string) InstrumentsServiceClient {
	conn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{AccessToken: token})),
	)
	if err != nil {
		panic(err)
	}

	return NewInstrumentsServiceClient(conn)
}

func CreateSandboxServiceClient(url string, token string) SandboxServiceClient {
	conn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{AccessToken: token})),
	)
	if err != nil {
		panic(err)
	}

	return NewSandboxServiceClient(conn)
}
