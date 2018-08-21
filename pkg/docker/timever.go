package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	TimeverRegexp = regexp.MustCompile(`^([0-9]{2})\.((?:0[1-9])|(?:1[0-2]))\.([0-9]+)(?:-(.+))?$`)

	ErrTimeverInvalidYear  = errors.New("invalid year value")
	ErrTimeverInvalidMonth = errors.New("invalid month value")
	ErrTimeverEmptyVersion = errors.New("version string empty")
	ErrTimeverParse        = errors.New("invalid Timever version string")
)

type Timever struct {
	Year    uint64
	Month   uint64
	Patch   uint64
	Release string
	// TODO: pre-release version support?
}

// Timever to string
func (v Timever) String() string {
	s := fmt.Sprintf("%02d.%02d.%d", v.Year, v.Month, v.Patch)
	if len(v.Release) > 0 {
		s += "-" + v.Release
	}
	return s
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (v Timever) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (v *Timever) UnmarshalJSON(data []byte) (err error) {
	var versionString string

	if err = json.Unmarshal(data, &versionString); err != nil {
		return
	}

	*v, err = TimeverParse(versionString)

	return
}

// EQ checks if v is equal to o.
func (v Timever) EQ(o Timever) bool {
	return (v.Compare(o) == 0)
}

// NE checks if v is not equal to o.
func (v Timever) NE(o Timever) bool {
	return (v.Compare(o) != 0)
}

// GT checks if v is greater than o.
func (v Timever) GT(o Timever) bool {
	return (v.Compare(o) == 1)
}

// GTE checks if v is greater than or equal to o.
func (v Timever) GTE(o Timever) bool {
	return (v.Compare(o) >= 0)
}

// GE checks if v is greater than or equal to o.
func (v Timever) GE(o Timever) bool {
	return (v.Compare(o) >= 0)
}

// LT checks if v is less than o.
func (v Timever) LT(o Timever) bool {
	return (v.Compare(o) == -1)
}

// LTE checks if v is less than or equal to o.
func (v Timever) LTE(o Timever) bool {
	return (v.Compare(o) <= 0)
}

// LE checks if v is less than or equal to o.
func (v Timever) LE(o Timever) bool {
	return (v.Compare(o) <= 0)
}

// Compare compares Timevers v to o:
// -1 == v is less than o
// 0 == v is equal to o
// 1 == v is greater than o
func (v Timever) Compare(o Timever) int {
	if v.Year != o.Year {
		if v.Year > o.Year {
			return 1
		}
		return -1
	}
	if v.Month != o.Month {
		if v.Month > o.Month {
			return 1
		}
		return -1
	}
	if v.Patch != o.Patch {
		if v.Patch > o.Patch {
			return 1
		}
		return -1
	}
	return 0
}

// Validate validates v and returns error in case
func (v Timever) Validate() error {
	if v.Year > 99 {
		return errors.New("Invalid year value")
	}
	if v.Month > 12 {
		return errors.New("Invalid month value")
	}
	return nil
}

// NewTimever is an alias for TimeverParse and returns a pointer, parses version string and returns a validated Timever or error
func NewTimever(s string) (vp *Timever, err error) {
	v, err := TimeverParse(s)
	vp = &v
	return
}

// TimeverParse parses version string and returns a validated Timever or error
func TimeverParse(s string) (Timever, error) {
	if len(s) == 0 {
		return Timever{}, errors.New("Version string empty")
	}
	match := TimeverRegexp.FindStringSubmatch(s)
	if len(match) != 5 {
		return Timever{}, errors.New("Invalid Timever version string")
	}
	year, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return Timever{}, err
	}
	month, err := strconv.ParseUint(match[2], 10, 64)
	if err != nil {
		return Timever{}, err
	}
	patch, err := strconv.ParseUint(match[3], 10, 64)
	if err != nil {
		return Timever{}, err
	}

	v := Timever{
		Year:    year,
		Month:   month,
		Patch:   patch,
		Release: match[4],
	}
	return v, nil
}

// MustTimeverParse is like Parse but panics if the version cannot be parsed.
func MustTimeverParse(s string) Timever {
	v, err := TimeverParse(s)
	if err != nil {
		panic(`timever: TimeverParse(` + s + `): ` + err.Error())
	}
	return v
}
