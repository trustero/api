package client

// Create common setup and teardown for dependent services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"strconv"
	"time"

	"github.com/trustero/api/go/config"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Ntrced *NtrcedClient

type NtrcedClient struct {
	Client        *grpc.ClientConn
	TlsDialOption grpc.DialOption
}

func InitGRPCClient(cert, override string) {
	Ntrced = &NtrcedClient{}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Err(err).Msg("error loading system cert pool")
		rootCAs = x509.NewCertPool()
	}

	if cert == "dev" {
		log.Info().Msgf("using dev cert with server name override %s", override)
		if !rootCAs.AppendCertsFromPEM([]byte(config.DevCertCa)) {
			panic("Error parsing X.509 certificate CA.")
		}
	}

	Ntrced.TlsDialOption = grpc.WithTransportCredentials(
		credentials.NewTLS(&tls.Config{
			ServerName: override,
			RootCAs:    rootCAs,
		}))
}

func (n *NtrcedClient) Dial(token, host string, port int) (err error) {
	// Dial options
	grpcCred := oauth.NewOauthAccess(&oauth2.Token{AccessToken: token})
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(grpcCred),
		n.TlsDialOption,
	}

	// Connect to local server
	addr := host + ":" + strconv.Itoa(port)
	n.Client, err = grpc.Dial(addr, opts...)
	return
}

func (n *NtrcedClient) CloseClient() error {
	if n.Client != nil {
		return n.Client.Close()
	}
	return nil
}

func (n *NtrcedClient) GetReceptorClient() (receptor_v1.ReceptorClient, context.Context, context.CancelFunc) {
	agt := receptor_v1.NewReceptorClient(n.Client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return agt, ctx, cancel
}
