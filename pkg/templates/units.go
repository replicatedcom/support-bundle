package templates

import units "github.com/docker/go-units"

func init() {
	RegisterFunc("bytesSize", units.BytesSize)
	RegisterFunc("fromHumanSize", units.FromHumanSize)
	RegisterFunc("humanDuration", units.HumanDuration)
	RegisterFunc("humanSize", units.HumanSize)
	RegisterFunc("humanSizeWithPrecision", units.HumanSizeWithPrecision)
	RegisterFunc("parseUlimit", units.ParseUlimit)
	RegisterFunc("ramInBytes", units.RAMInBytes)
}
