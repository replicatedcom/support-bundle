package docker

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/blang/semver"
)

var (
	DockerReleaseCE = "ce"
	DockerReleaseEE = "ee"

	ErrServerVersionInvalidRelease = errors.New("release must be one of ce or ee")
)

type ServerVersion struct {
	Semver  *semver.Version
	Timever *Timever
}

func (v ServerVersion) Release() string {
	if v.Timever == nil {
		return ""
	}
	return strings.SplitN(v.Timever.Release, "-", 2)[0]
}

func (v ServerVersion) String() string {
	if v.Timever != nil {
		return v.Timever.String()
	}
	if v.Semver != nil {
		return v.Semver.String()
	}
	return ""
}

func (v ServerVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *ServerVersion) UnmarshalJSON(data []byte) (err error) {
	var versionString string

	if err = json.Unmarshal(data, &versionString); err != nil {
		return
	}

	*v, err = ParseServerVersion(versionString)
	return
}

func (v ServerVersion) EQ(o ServerVersion) bool {
	return (v.Compare(o) == 0)
}

func (v ServerVersion) NE(o ServerVersion) bool {
	return (v.Compare(o) != 0)
}

func (v ServerVersion) GT(o ServerVersion) bool {
	return (v.Compare(o) == 1)
}

func (v ServerVersion) GTE(o ServerVersion) bool {
	return (v.Compare(o) >= 0)
}

func (v ServerVersion) GE(o ServerVersion) bool {
	return (v.Compare(o) >= 0)
}

func (v ServerVersion) LT(o ServerVersion) bool {
	return (v.Compare(o) == -1)
}

func (v ServerVersion) LTE(o ServerVersion) bool {
	return (v.Compare(o) <= 0)
}

func (v ServerVersion) LE(o ServerVersion) bool {
	return (v.Compare(o) <= 0)
}

func (v ServerVersion) Compare(o ServerVersion) int {
	if v.Timever != nil && o.Timever != nil {
		return v.Timever.Compare(*o.Timever)
	}
	if v.Semver != nil && o.Semver != nil {
		return v.Semver.Compare(*o.Semver)
	}
	if v.Timever != nil {
		return 1
	}
	if o.Timever != nil {
		return -1
	}
	if v.Semver != nil {
		return 1
	}
	if o.Semver != nil {
		return -1
	}
	return 0
}

func (v ServerVersion) Validate() error {
	if v.Timever != nil {
		if err := v.Timever.Validate(); err != nil {
			return err
		}
		return v.ValidateRelease()
	}
	if v.Semver != nil {
		return v.Semver.Validate()
	}
	return nil
}

func (v ServerVersion) ValidateRelease() error {
	release := v.Release()
	if release != "" && release != DockerReleaseCE && release != DockerReleaseEE {
		return ErrServerVersionInvalidRelease
	}
	return nil
}

func NewServerVersion(s string) (vp *ServerVersion, err error) {
	v, err := ParseServerVersion(s)
	vp = &v
	return
}

func ParseServerVersion(s string) (ServerVersion, error) {
	if s == "" {
		return ServerVersion{}, nil
	}
	timever, timeverErr := TimeverParse(s)
	if timeverErr == nil {
		if timever.Year >= 17 {
			v := ServerVersion{Timever: &timever}
			if err := v.ValidateRelease(); err != nil {
				return ServerVersion{}, err
			}
			return v, nil
		}
	}
	ver, semverErr := semver.Parse(s)
	if semverErr == nil {
		return ServerVersion{Semver: &ver}, nil
	}
	if timeverErr != nil {
		return ServerVersion{}, timeverErr
	}
	return ServerVersion{}, semverErr
}

func MustParseServerVersion(s string) ServerVersion {
	v, err := ParseServerVersion(s)
	if err != nil {
		panic(`ParseServerVersion(` + s + `): ` + err.Error())
	}
	return v
}
