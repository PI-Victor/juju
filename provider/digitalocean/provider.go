package digitalocean

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/juju/jsonschema"
	"github.com/juju/loggo"

	"github.com/digitalocean/godo"
	"github.com/juju/errors"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	envCtx "github.com/juju/juju/environs/context"
)

var logger = loggo.GetLogger("juju.provider.digitalocean")

type tokenSource struct {
	AccessToken string
}

// NOTE: done to satisfy the token interface for oauth2.
func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type environProvider struct {
	environProviderCredentials
}

var providerInstance environProvider

func (environProvider) Version() int {
	return 0
}

func (p environProvider) Open(args environs.OpenParams) (environs.Environ, error) {
	logger.Infof("opening model %q", args.Config.Name())

	e := new(environ)
	e.cloud = args.Cloud
	e.name = args.Config.Name()
	e.dgo = newDGOClient(e.cloud)
	return e, nil
}

func newDGOClient(cloud environs.CloudSpec) *godo.Client {
	creds := cloud.Credential.Attributes()
	t := &tokenSource{
		AccessToken: creds["AccessToken"],
	}
	newOauth2 := oauth2.NewClient(context.Background(), t)
	return godo.NewClient(newOauth2)
}

func (p environProvider) CloudSchema() *jsonschema.Schema {
	return nil
}

func newDOClient(cloud environs.CloudSpec) (*godo.Client, error) {
	return nil, nil
}

func (environProvider) newConfig(cfg *config.Config) error {
	return nil
}

// Ping tests the connection to the cloud, to verify the endpoint is valid.
func (p environProvider) Ping(ctx envCtx.ProviderCallContext, endpoint string) error {
	return errors.NotImplementedf("Ping")
}
