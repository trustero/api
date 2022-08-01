package client

// Create common setup and teardown for dependent services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"strconv"
	"time"

	"github.com/trustero/api/go/receptor_sdk/config"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var ServerConn *ServerConnection

type ServerConnection struct {
	Connection    *grpc.ClientConn
	TlsDialOption grpc.DialOption
}

func InitGRPCClient(cert, override string) {
	ServerConn = &ServerConnection{}

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

	ServerConn.TlsDialOption = grpc.WithTransportCredentials(
		credentials.NewTLS(&tls.Config{
			ServerName: override,
			RootCAs:    rootCAs,
		}))
}

func (sc *ServerConnection) Dial(token, host string, port int) (err error) {
	// Dial options
	grpcCred := oauth.NewOauthAccess(&oauth2.Token{AccessToken: token})
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(logUnaryCall),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(grpcCred),
		sc.TlsDialOption,
	}

	// Connect to local server
	addr := host + ":" + strconv.Itoa(port)
	sc.Connection, err = grpc.Dial(addr, opts...)
	return
}

func (sc *ServerConnection) CloseClient() error {
	if sc.Connection != nil {
		return sc.Connection.Close()
	}
	return nil
}

func (sc *ServerConnection) GetReceptorClient() (rc receptor_v1.ReceptorClient) {
	rc = receptor_v1.NewReceptorClient(sc.Connection)
	return
}

func logUnaryCall(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	callId := RandString(8)
	log.Trace().Msgf("grpc begin [%s-%s]", method, callId)
	defer func() {
		elapsed := time.Now().Sub(start)
		log.Info().Msgf("grpc end [%s-%s], elapsed time: %fs", method, callId, elapsed.Seconds())
	}()

	return invoker(ctx, method, req, reply, cc, opts...)
}
