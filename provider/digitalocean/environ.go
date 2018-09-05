package digitalocean

import (
	"sync"



	"github.com/digitalocean/godo"

	"github.com/juju/errors"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	envCtx "github.com/juju/juju/environs/context"
	"github.com/juju/juju/instance"
	"github.com/juju/version"
)

type environ struct {
	name string

	cloud environs.CloudSpec
	dgo   *godo.Client

	lock sync.Mutex // lock protects access to ecfg
	ecfg *environConfig
}

var (
	_ environs.Environ           = (*environ)(nil)
	_ environs.NetworkingEnviron = (*environ)(nil)
)

func (e *environ) Name() string {
	return e.name
}

func (e *environ) Create(envCtx.ProviderCallContext, environs.CreateParams) error {
	return nil
}

func (e *environ) SetConfig() error {
	cfg, err := providerInstance.newConfig(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}
func (e environ) AdoptResources(ctx envCtx.ProviderCallContext, controllerUUID string, fromVersion version.Number) error {
	return nil
}

func (e *environ) Region() string {
	return ""
}

func (e *environ) Provider() environs.EnvironProvider {
	return providerInstance
}

func (e *environ) Config() *config.Config {
	return e.ecfg.config
}

func (e *environ) PrepareForBootstrap(ctx environs.BootstrapContext) error {
	return nil
}

func (e *environ) Bootstrap(ctx environs.BootstrapContext, callCtx envCtx.ProviderCallContext, params environs.BootstrapParams) (*environs.BootstrapResult, error) {
	return nil, nil
}

func (e *environ) AllInstances(ctx envCtx.ProviderCallContext) ([]instance.Instance, error) {
	return []instance.Instance{}, nil
}

func (e *environ) ConstraintsValidator(ctx envCtx.ProviderCallContext) (constraints.Validator, error) {
	validator := constraints.NewValidator()
	return validator, nil
}

func (e *environ) ControllerInstances(ctx envCtx.ProviderCallContext, controllerUUID string) ([]instance.Id, error) {
	return []instance.Id{}, nil
}

func (e *environ) Destroy(ctx envCtx.ProviderCallContext) error {
	return nil
}

func (e *environ) DestroyController(ctx envCtx.ProviderCallContext, controllerUUID string) error {
	return nil
}

func (e *environ) supportedInstanceTypes(ctx envCtx.ProviderCallContext) ([]instances.InstanceType, error) {
	return nil
}
