// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Package client provides GRPC utilities to simplify connecting to Trustero GRPC service.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"strconv"
	"time"

	"github.com/trustero/api/go/receptor_sdk/config"
	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	grpcCred := oauth.TokenSource{TokenSource: ts}
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(logUnaryCall),
		grpc.WithPerRPCCredentials(grpcCred),
		sc.TlsDialOption,
		grpc.WithStreamInterceptor(logStreamCall),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(500 * 1024 * 1024)),
	}

	// Connect to local server
	addr := host + ":" + strconv.Itoa(port)
	sc.Connection, err = grpc.NewClient(addr, opts...)
	return
}

// DialAndWait dials and waits for the connection to reach Ready state or timeout.
// This preserves old blocking behavior safely, avoiding deprecated WithBlock.
func (sc *ServerConnection) DialAndWait(token, host string, port int, timeout time.Duration) error {
	if err := sc.Dial(token, host, port); err != nil {
		return err
	}
	if sc.Connection == nil {
		return errors.New("grpc connection is nil after dial")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		state := sc.Connection.GetState()
		if state == connectivity.Ready {
			return nil
		}
		sc.Connection.Connect()
		if ok := sc.Connection.WaitForStateChange(ctx, state); !ok {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return errors.New("grpc wait for state change failed")
		}
	}
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
