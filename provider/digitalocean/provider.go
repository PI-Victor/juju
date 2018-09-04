package digitalocean

import (
	"github.com/juju/loggo"

	"github.com/digitalocean/godo"
	"github.com/juju/juju/environs"
)

var logger = loggo.GetLogger("juju.provider.digitalocean")

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

}

func digitalOceanClient(cloud environs.CloudSpec) (*godo.Client, error) {

	return nil, nil
}
