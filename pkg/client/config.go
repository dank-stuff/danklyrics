package client

import "codeberg.org/dankstuff/danklyrics/pkg/provider"

// Config holds the configs needed to initialize [Local] or [Http] clients.
type Config struct {
	Providers     []provider.Name
	ProvidersAuth map[provider.Name]provider.Auth
	// ApiAddress only used by [Http] client, setting its value for [Local] client won't destroy the world, but it's pointless.
	// defaults to (https://api.danklyrics.com)
	ApiAddress string
}
