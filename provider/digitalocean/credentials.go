package digitalocean

import (
	"github.com/juju/juju/cloud"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/context"
)

type environProviderCredentials struct{}

func (environProviderCredentials) CredentialSchemas() map[cloud.AuthType]cloud.CredentialSchema {
	return map[cloud.AuthType]cloud.CredentialSchema{}
}

func (environProviderCredentials) DetectCredentials() (*cloud.CloudCredential, error) {
	return nil, nil
}

func (environProviderCredentials) FinalizeCredential(_ environs.FinalizeCredentialContext, _ environs.FinalizeCredentialParams) (*cloud.Credential, error) {
	return nil, nil
}

func (environProviderCredentials) Ping(ctx context.ProviderCallContext, endpoint string) error {
	return nil
}

func (environProviderCredentials) PrepareConfig(environs.PrepareConfigParams) (*config.Config, error) {
	return nil, nil
}
