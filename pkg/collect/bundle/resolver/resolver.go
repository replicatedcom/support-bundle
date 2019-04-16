package resolver

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var (
	DefaultResolvers = []Resolver{
		new(S3Resolver),
		new(FileResolver),
	}

	resolversMap = map[string]Resolver{
		"file": new(FileResolver),
		"s3":   new(S3Resolver),
	}
)

type Factory struct {
	Logger    log.Logger
	Resolvers []Resolver
}

type Resolver interface {
	Detect(name string) (string, bool, error)
	Resolve(name string) (afero.Fs, string, error)
}

func NewWithDefaults(logger log.Logger) *Factory {
	return &Factory{
		Logger:    logger,
		Resolvers: DefaultResolvers,
	}
}

func New(logger log.Logger, resolvers []Resolver) *Factory {
	if resolvers == nil {
		return NewWithDefaults(logger)
	}
	return &Factory{
		Logger:    logger,
		Resolvers: resolvers,
	}
}

func (f *Factory) Fs(name string) (afero.Fs, string, error) {
	forced, name := getForcedGetter(name)
	if forced != "" {
		resolver, ok := resolversMap[forced]
		if !ok {
			return nil, "", fmt.Errorf("unknown forced resolver %s", forced)
		}
		fs, path, err := resolver.Resolve(name)
		return fs, path, errors.Wrapf(err, "resolve %T", resolver)
	}
	for _, resolver := range f.Resolvers {
		name, ok, err := resolver.Detect(name)
		if err != nil {
			return nil, "", errors.Wrapf(err, "detect %T", resolver)
		} else if !ok {
			continue
		}
		fs, path, err := resolver.Resolve(name)
		return fs, path, errors.Wrapf(err, "resolve %T", resolver)
	}
	return nil, "", errors.New("unsupported fs")
}
