package resolver

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var (
	DefaultResolvers = []Resolver{
		new(S3Resolver),
		new(FileResolver),
	}
)

type Factory struct {
	resolvers []Resolver
}

type Resolver interface {
	Detect(name string) (string, bool, error)
	Resolve(name string) (afero.Fs, string, error)
}

func NewWithDefaults() *Factory {
	return &Factory{
		resolvers: DefaultResolvers,
	}
}

func New(resolvers []Resolver) *Factory {
	if resolvers == nil {
		return NewWithDefaults()
	}
	return &Factory{
		resolvers: resolvers,
	}
}

func (f *Factory) Fs(name string) (afero.Fs, string, error) {
	for _, resolver := range f.resolvers {
		forced, ok, err := resolver.Detect(name)
		if err != nil {
			return nil, "", errors.Wrapf(err, "detect %T", resolver)
		} else if !ok {
			continue
		}
		fs, path, err := resolver.Resolve(forced)
		return fs, path, errors.Wrapf(err, "resolve %T", resolver)
	}
	return nil, "", errors.New("unsupported fs")
}
