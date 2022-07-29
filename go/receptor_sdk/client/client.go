package client

// Create common setup and teardown for dependent services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/trustero/api/go/receptor_sdk/config"
	"google.golang.org/grpc/credentials"
	"strconv"
	"time"

	"github.com/trustero/api/go/receptor_v1"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type AuthClientCallback func(context.Context, receptor_v1.ReceptorClient, *receptor_v1.ReceptorConfiguration) error

type Server struct {
	Host               string
	Port               int
	Cert               string
	CertServerOverride string
	Token              string
}

func (s *Server) toString() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

type ScopedClientFactory interface {
	AuthScope(receptorModelId string, onClient AuthClientCallback) error
}

type Context struct {
	Server               *Server
	connection           *grpc.ClientConn
	TransportCredentials credentials.TransportCredentials
}

var Factory ScopedClientFactory

func InitFactory(server *Server) ScopedClientFactory {
	var rootCAs *x509.CertPool
	var err error
	if rootCAs, err = x509.SystemCertPool(); err != nil {
		log.Warn().Msg("error loading system cert pool. Assigning and empty cert pool")
		rootCAs = x509.NewCertPool()
		err = nil
	}

	if server.Cert == "dev" {
		log.Info().Msgf("using dev cert with server name override %s", server.CertServerOverride)
		if !rootCAs.AppendCertsFromPEM([]byte(config.DevCertCa)) {
			panic("Error parsing X.509 certificate CA.")
		}
	}

	return &Context{
		Server: server,
		TransportCredentials: credentials.NewTLS(&tls.Config{
			ServerName: server.CertServerOverride,
			RootCAs:    rootCAs,
		}),
	}
}

func (t *Context) openConnection() (err error) {
	// Dial options
	grpcCred := oauth.NewOauthAccess(&oauth2.Token{AccessToken: t.Server.Token})
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(grpcCred),
		grpc.WithTransportCredentials(t.TransportCredentials),
	}
	t.connection, err = grpc.Dial(t.Server.toString(), opts...)
	return
}

func (t *Context) closeConnection() error {
	if t.connection != nil {
		return t.connection.Close()
	}
	return nil
}

func (t *Context) AuthScope(modelId string, onClient AuthClientCallback) (err error) {
	if err = t.openConnection(); err != nil {
		return
	}
	defer func() {
		if err = t.closeConnection(); err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	client := receptor_v1.NewReceptorClient(t.connection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get credentials and config
	var receptorConfiguration *receptor_v1.ReceptorConfiguration
	// FIXME: The receptor doesnt know it's receptor object ID at this point.
	if receptorConfiguration, err = client.GetConfiguration(ctx, &receptor_v1.ReceptorOID{ReceptorObjectId: modelId}); err != nil {
		return err
	}
	return onClient(ctx, client, receptorConfiguration)
}
