package digitalocean

import (
	"github.com/juju/juju/environs"
)

const (
	providerType = "digitalocean"
)

func init() {
	environs.RegisterProvider(providerType, providerInstance)
}
