// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package client provides GRPC utilities to simplify connecting to Trustero GRPC service.
package client

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

// InitGRPCClient sets up the SSL certificate for subsequent Trustero GRPC connections.
func InitGRPCClient(cert, override string) {
	ServerConn = &ServerConnection{}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Err(err).Msg("error loading system certificate pool")
		rootCAs = x509.NewCertPool()
	}

	if cert == "dev" {
		log.Info().Msgf("using a development SSL certificate with server name override %s", override)
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

// Dial makes a GRPC connection to Trustero GRPC service.  A Trustero JWT bearer token must be provided.
func (sc *ServerConnection) Dial(token, host string, port int) (err error) {
	// Dial options
	grpcCred := oauth.NewOauthAccess(&oauth2.Token{AccessToken: token})
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(logUnaryCall),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(grpcCred),
		sc.TlsDialOption,
		grpc.WithStreamInterceptor(logStreamCall),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(500 * 1024 * 1024)),
	}

	// Connect to local server
	addr := host + ":" + strconv.Itoa(port)
	sc.Connection, err = grpc.Dial(addr, opts...)
	return
}

// CloseClient closes a previously dialed Trustero GRPC connection.
func (sc *ServerConnection) CloseClient() error {
	if sc.Connection != nil {
		return sc.Connection.Close()
	}
	return nil
}

// GetReceptorClient returns a Go client instance that implements the [receptor_v1.Receptor] Protobuf interface.
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

func logStreamCall(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	start := time.Now()
	callId := RandString(8)
	log.Trace().Msgf("grpc begin stream [%s-%s]", method, callId)
	defer func() {
		elapsed := time.Now().Sub(start)
		log.Info().Msgf("grpc end stream [%s-%s], elapsed time: %fs", method, callId, elapsed.Seconds())
	}()

	// Invoke the streamer and return the stream
	clientStream, err := streamer(ctx, desc, cc, method, opts...)
	return clientStream, err
}
