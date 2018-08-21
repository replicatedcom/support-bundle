package docker

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var timeverFormatTestCases = []struct {
	v      Timever
	result string
}{
	{Timever{Year: 7, Month: 3, Patch: 123, Release: ""}, "07.03.123"},
	{Timever{Year: 17, Month: 3, Patch: 123, Release: "ce"}, "17.03.123-ce"},
	{Timever{Year: 17, Month: 3, Patch: 0, Release: "ee-1"}, "17.03.0-ee-1"},
}

func TestTimeverStringer(t *testing.T) {
	for _, test := range timeverFormatTestCases {
		if res := test.v.String(); res != test.result {
			t.Errorf("Stringer, expected %q but got %q", test.result, res)
		}
	}
}

func TestTimeverParse(t *testing.T) {
	for _, test := range timeverFormatTestCases {
		if v, err := TimeverParse(test.result); err != nil {
			t.Errorf("Error parsing %q: %q", test.result, err)
		} else if comp := v.Compare(test.v); comp != 0 {
			t.Errorf("Parsing, expected %q but got %q, comp: %d ", test.v, v, comp)
		} else if err := v.Validate(); err != nil {
			t.Errorf("Error validating parsed version %q: %q", v, err)
		}
	}
}

func TestMustTimeverParse(t *testing.T) {
	_ = MustTimeverParse("17.03.123-ce")
}

func TestMustTimeverParse_panic(t *testing.T) {
	assert.Panics(t, func() {
		MustTimeverParse("invalid version")
	})
}

var timeverCompareTestCases = []struct {
	v1     Timever
	v2     Timever
	result int
}{
	{Timever{Year: 1, Month: 0, Patch: 0, Release: ""}, Timever{Year: 1, Month: 0, Patch: 0, Release: ""}, 0},
	{Timever{Year: 2, Month: 0, Patch: 0, Release: ""}, Timever{Year: 1, Month: 0, Patch: 0, Release: ""}, 1},
	{Timever{Year: 0, Month: 1, Patch: 0, Release: ""}, Timever{Year: 0, Month: 1, Patch: 0, Release: ""}, 0},
	{Timever{Year: 0, Month: 2, Patch: 0, Release: ""}, Timever{Year: 0, Month: 1, Patch: 0, Release: ""}, 1},
	{Timever{Year: 0, Month: 0, Patch: 1, Release: ""}, Timever{Year: 0, Month: 0, Patch: 1, Release: ""}, 0},
	{Timever{Year: 0, Month: 0, Patch: 2, Release: ""}, Timever{Year: 0, Month: 0, Patch: 1, Release: ""}, 1},
	{Timever{Year: 1, Month: 2, Patch: 3, Release: ""}, Timever{Year: 1, Month: 2, Patch: 3, Release: ""}, 0},
	{Timever{Year: 2, Month: 2, Patch: 4, Release: ""}, Timever{Year: 1, Month: 2, Patch: 4, Release: ""}, 1},
	{Timever{Year: 1, Month: 3, Patch: 3, Release: ""}, Timever{Year: 1, Month: 2, Patch: 3, Release: ""}, 1},
	{Timever{Year: 1, Month: 2, Patch: 4, Release: ""}, Timever{Year: 1, Month: 2, Patch: 3, Release: ""}, 1},
	{Timever{Year: 1, Month: 0, Patch: 0, Release: ""}, Timever{Year: 2, Month: 0, Patch: 0, Release: ""}, -1},
	{Timever{Year: 2, Month: 0, Patch: 0, Release: ""}, Timever{Year: 2, Month: 1, Patch: 0, Release: ""}, -1},
	{Timever{Year: 2, Month: 1, Patch: 0, Release: ""}, Timever{Year: 2, Month: 1, Patch: 1, Release: ""}, -1},
	{Timever{Year: 1, Month: 0, Patch: 0, Release: "ce"}, Timever{Year: 1, Month: 0, Patch: 0, Release: "ee"}, 0},
}

func TestTimeverCompare(t *testing.T) {
	for _, test := range timeverCompareTestCases {
		if res := test.v1.Compare(test.v2); res != test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v1, test.v2, test.result, res)
		}
		//Test counterpart
		if res := test.v2.Compare(test.v1); res != -test.result {
			t.Errorf("Comparing %q : %q, expected %d but got %d", test.v2, test.v1, -test.result, res)
		}
	}
}

func TestTimeverCompareHelper(t *testing.T) {
	v := Timever{1, 0, 0, "ce"}
	v1 := Timever{1, 0, 1, "ce"}
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

var wrongTimeverFormatTestCases = []struct {
	str string
}{
	{str: ""},
	{str: "."},
	{str: "1."},
	{str: ".1"},
	{str: "a.b.c"},
	{str: "1.a.b"},
	{str: "1.1.a"},
	{str: "1.a.1"},
	{str: "a.1.1"},
	{str: ".."},
	{str: "1.."},
	{str: "1.1."},
	{str: "1..1"},
	{str: "-1.1.1"},
	{str: "1.-1.1"},
	{str: "1.1.-1"},
	{str: "100.1.1"},
	{str: "1.13.1"},
}

func TestTimeverWrongFormat(t *testing.T) {
	for _, test := range wrongTimeverFormatTestCases {
		if res, err := TimeverParse(test.str); err == nil {
			t.Errorf("Parsing wrong format version %q, expected error but got %q", test.str, res)
		}
	}
}

func TestTimeverJSONMarshal(t *testing.T) {
	versionString := "17.03.123-ce"
	v, err := TimeverParse(versionString)
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
}

func TestTimeverJSONUnmarshal(t *testing.T) {
	versionString := "17.03.123-ce"
	quotedVersionString := strconv.Quote(versionString)

	var v Timever
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
