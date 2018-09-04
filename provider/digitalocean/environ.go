package digitalocean

import (
	"context"
	"sync"

	"golang.org/x/oauth2"

	"github.com/digitalocean/godo"

	"github.com/juju/errors"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/instance"
)

type doConnection interface {
	VerifyCredentials() error

	Instance(id, zone string) (godo.Droplet, error)
	Instances(prefix string, statuses ...string) ([]godo.Droplet, error)
}

type environ struct {
	name  string
	uuid  string
	cloud environs.CloudSpec
	dgo   doConnection

	lock sync.Mutex // lock protects access to ecfg
	ecfg *environConfig

	namespace instance.Namespace
}

var (
	_ environs.Environ           = (*environ)(nil)
	_ environs.NetworkingEnviron = (*environ)(nil)
)

func (e *environ) Name() string {
	return e.name
}

func (e *environ) SetConfig() string {
	return
}

func (e *environ) Region() string {
	return
}

func (e *environ) Provider() environs.EnvironProvider {
	return providerInstance
}

func (e *environ) Config() *config.Config {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.ecfg.config
}

func (e *environ) PrepareForBootstrap(ctx environs.BootstrapContext) error {
	if ctx.ShouldVerifyCredentials() {
		if err := e.dgo.VerifyCredentials(); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}
