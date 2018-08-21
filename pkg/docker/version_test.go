package docker

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
)

var serverVersionFormatTestCases = []struct {
	v      ServerVersion
	result string
}{
	{ServerVersion{Timever: &Timever{Year: 17, Month: 2, Patch: 3, Release: ""}}, "17.02.3"},
	{ServerVersion{Timever: &Timever{Year: 17, Month: 2, Patch: 3, Release: "ce"}}, "17.02.3-ce"},
	{ServerVersion{Timever: &Timever{Year: 17, Month: 2, Patch: 3, Release: "ee-1"}}, "17.02.3-ee-1"},
	{ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}, "1.2.3"},
	{ServerVersion{}, ""},
}

func TestServerVersionStringer(t *testing.T) {
	for _, test := range serverVersionFormatTestCases {
		if res := test.v.String(); res != test.result {
			t.Errorf("Stringer, expected %q but got %q", test.result, res)
		}
	}
}

func TestServerVersionParse(t *testing.T) {
	for _, test := range serverVersionFormatTestCases {
		if v, err := ParseServerVersion(test.result); err != nil {
			t.Errorf("Error parsing %q: %q", test.result, err)
		} else if comp := v.Compare(test.v); comp != 0 {
			t.Errorf("Parsing, expected %q but got %q, comp: %d ", test.v, v, comp)
		} else if err := v.Validate(); err != nil {
			t.Errorf("Error validating parsed version %q: %q", v, err)
		}
	}
}

func TestMustServerVersionParse(t *testing.T) {
	_ = MustParseServerVersion("17.03.123-ce")
}

func TestMustServerVersionParse_panic(t *testing.T) {
	assert.Panics(t, func() {
		MustParseServerVersion("invalid version")
	})
}

var serverVersionCompareTestCases = []struct {
	v1     ServerVersion
	v2     ServerVersion
	result int
}{
	{ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: ""}}, ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: ""}}, 0},
	{ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: "ce"}}, ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: "ee"}}, 0},
	{ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}, ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}, 0},
	{ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: ""}}, ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}, 1},
	{ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}, ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: ""}}, -1},
	{ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: ""}}, ServerVersion{}, 1},
	{ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}, ServerVersion{}, 1},
	{ServerVersion{}, ServerVersion{}, 0},
}

func TestServerVersionCompare(t *testing.T) {
	for _, test := range serverVersionCompareTestCases {
		if res := test.v1.Compare(test.v2); res != test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v1, test.v2, test.result, res)
		}
		//Test counterpart
		if res := test.v2.Compare(test.v1); res != -test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v2, test.v1, -test.result, res)
		}
	}
}

func TestServerVersionCompareHelper(t *testing.T) {
	v := ServerVersion{Semver: &semver.Version{Major: 1, Minor: 2, Patch: 3, Pre: nil, Build: nil}}
	v1 := ServerVersion{Timever: &Timever{Year: 1, Month: 2, Patch: 3, Release: ""}}
	if !v.EQ(v) {
		t.Errorf("%q should be equal to %q", v, v)
	}
	if !v1.NE(v) {
		t.Errorf("%q should not be equal to %q", v1, v)
	}
	if !v.GTE(v) {
		t.Errorf("%q should be greater than or equal to %q", v, v)
	}
	if !v.LTE(v) {
		t.Errorf("%q should be less than or equal to %q", v, v)
	}
	if !v.LT(v1) {
		t.Errorf("%q should be less than %q", v, v1)
	}
	if !v.LTE(v1) {
		t.Errorf("%q should be less than or equal %q", v, v1)
	}
	if !v.LE(v1) {
		t.Errorf("%q should be less than or equal %q", v, v1)
	}
	if !v1.GT(v) {
		t.Errorf("%q should be greater than %q", v1, v)
	}
	if !v1.GTE(v) {
		t.Errorf("%q should be greater than or equal %q", v1, v)
	}
	if !v1.GE(v) {
		t.Errorf("%q should be greater than or equal %q", v1, v)
	}
}

type serverVersionWrongformatTest struct {
	str string
}

var serverVersionWrongformatTests = []serverVersionWrongformatTest{
	{"a.b.c"},
	{"17.03.0-invalid"},
}

func TestServerVersionWrongFormat(t *testing.T) {
	for _, test := range serverVersionWrongformatTests {
		if res, err := ParseServerVersion(test.str); err == nil {
			t.Errorf("Parsing wrong format version %q, expected error but got %q", test.str, res)
		}
	}
}

func TestServerVersionJSONMarshal(t *testing.T) {
	versionString := "17.03.123-ce"
	v, err := ParseServerVersion(versionString)
	if err != nil {
		t.Error(err)
	}

	versionJSON, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	quotedVersionString := strconv.Quote(versionString)

	if string(versionJSON) != quotedVersionString {
		t.Errorf("JSON marshaled semantic version not equal: expected %q, got %q", quotedVersionString, string(versionJSON))
	}

	versionString = "1.13.1"
	v, err = ParseServerVersion(versionString)
	if err != nil {
		t.Error(err)
	}

	versionJSON, err = json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	quotedVersionString = strconv.Quote(versionString)

	if string(versionJSON) != quotedVersionString {
		t.Errorf("JSON marshaled semantic version not equal: expected %q, got %q", quotedVersionString, string(versionJSON))
	}
}

func TestServerVersionJSONUnmarshal(t *testing.T) {
	versionString := "17.03.123-ce"
	quotedVersionString := strconv.Quote(versionString)

	var v ServerVersion
	if err := json.Unmarshal([]byte(quotedVersionString), &v); err != nil {
		t.Error(err)
	}

	if v.String() != versionString {
		t.Errorf("JSON unmarshaled semantic version not equal: expected %q, got %q", versionString, v.String())
	}

	versionString = "1.13.1"
	quotedVersionString = strconv.Quote(versionString)

	if err := json.Unmarshal([]byte(quotedVersionString), &v); err != nil {
		t.Error(err)
	}

	if v.String() != versionString {
		t.Errorf("JSON unmarshaled semantic version not equal: expected %q, got %q", versionString, v.String())
	}

	badVersionString := strconv.Quote("a.b.c")
	if err := json.Unmarshal([]byte(badVersionString), &v); err == nil {
		t.Error("expected JSON unmarshal error, got nil")
	}
}
