package digitalocean

import (
	"github.com/juju/juju/environs"
)

const (
	provider = "digitalocean"
)

func init() {
	environs.RegisterProvider(providerType, providerInstance)
}
