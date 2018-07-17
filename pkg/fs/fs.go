package fs

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func FromViper(v *viper.Viper) afero.Fs {
	return New()
}

func New() afero.Fs {
	return &afero.Afero{Fs: afero.NewOsFs()}
}
