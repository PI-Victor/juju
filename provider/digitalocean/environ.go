package digitalocean

import (
	"sync"

	"github.com/digitalocean/godo"

	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	envCtx "github.com/juju/juju/environs/context"
	"github.com/juju/juju/environs/instances"
	"github.com/juju/juju/instance"
	"github.com/juju/juju/storage"
	"github.com/juju/version"
)

type doEnviron struct {
	name string

	cloud environs.CloudSpec
	dgo   *godo.Client

	lock sync.Mutex // lock protects access to ecfg
	ecfg *environConfig
}

var (
	_ environs.Environ = (*doEnviron)(nil)
)

func (e *doEnviron) Name() string {
	return e.name
}

func (e *doEnviron) Create(envCtx.ProviderCallContext, environs.CreateParams) error {
	return nil
}

func (e *doEnviron) SetConfig(cfg *config.Config) error {
	return providerInstance.newConfig(cfg)
}
func (e *doEnviron) AdoptResources(
	ctx envCtx.ProviderCallContext,
	controllerUUID string,
	fromVersion version.Number,
) error {
	return nil
}

func (e *doEnviron) Region() string {
	return ""
}

func (e *doEnviron) Provider() environs.EnvironProvider {
	return providerInstance
}

func (e *doEnviron) Config() *config.Config {
	return e.ecfg.config
}

func (e *doEnviron) PrepareForBootstrap(ctx environs.BootstrapContext) error {
	return nil
}

func (e *doEnviron) Bootstrap(
	ctx environs.BootstrapContext,
	callCtx envCtx.ProviderCallContext,
	params environs.BootstrapParams,
) (*environs.BootstrapResult, error) {
	return nil, nil
}

func (e *doEnviron) AllInstances(ctx envCtx.ProviderCallContext) ([]instance.Instance, error) {
	return []instance.Instance{}, nil
}

func (e *doEnviron) StorageProvider(t storage.ProviderType) (storage.Provider, error) {
	return &doStorageProvider{}, nil
}

func (e *doEnviron) StorageProviderTypes() ([]storage.ProviderType, error) {
	return nil, nil
}

func (e *doEnviron) ConstraintsValidator(ctx envCtx.ProviderCallContext) (constraints.Validator, error) {
	validator := constraints.NewValidator()
	return validator, nil
}

func (e *doEnviron) ControllerInstances(ctx envCtx.ProviderCallContext, controllerUUID string) ([]instance.Id, error) {
	return []instance.Id{}, nil
}

func (e *doEnviron) Destroy(ctx envCtx.ProviderCallContext) error {
	return nil
}

func (e *doEnviron) DestroyController(ctx envCtx.ProviderCallContext, controllerUUID string) error {
	return nil
}

func (e *doEnviron) Instances(ctx envCtx.ProviderCallContext, ids []instance.Id) ([]instance.Instance, error) {
	return []instance.Instance{}, nil
}

func (e *doEnviron) InstanceTypes(envCtx.ProviderCallContext, constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
	return instances.InstanceTypesWithCostMetadata{}, nil
}

func (e *doEnviron) MaintainInstance(envCtx.ProviderCallContext, environs.StartInstanceParams) error {
	return nil
}

func (e *doEnviron) PrecheckInstance(ctx envCtx.ProviderCallContext, args environs.PrecheckInstanceParams) error {
	return nil
}

func (e *doEnviron) StartInstance(ctx envCtx.ProviderCallContext, args environs.StartInstanceParams) (*environs.StartInstanceResult, error) {
	return nil, nil
}

func (e *doEnviron) StopInstances(ctx envCtx.ProviderCallContext, ids ...instance.Id) error {
	return nil
}
