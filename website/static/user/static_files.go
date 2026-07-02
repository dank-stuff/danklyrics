package staticuser

import "embed"

//go:embed *
var publicFiles embed.FS

func FS() embed.FS {
	return publicFiles
}
