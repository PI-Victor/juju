package digitalocean

import (
	"github.com/juju/juju/environs/config"
)

type environConfig struct {
	config *config.Config
	attrs  map[string]interface{}
}
