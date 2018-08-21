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
	{Timever{7, 3, 123, ""}, "07.03.123"},
	{Timever{17, 3, 123, "ce"}, "17.03.123-ce"},
	{Timever{17, 3, 0, "ee-1"}, "17.03.0-ee-1"},
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
	{Timever{1, 0, 0, ""}, Timever{1, 0, 0, ""}, 0},
	{Timever{2, 0, 0, ""}, Timever{1, 0, 0, ""}, 1},
	{Timever{0, 1, 0, ""}, Timever{0, 1, 0, ""}, 0},
	{Timever{0, 2, 0, ""}, Timever{0, 1, 0, ""}, 1},
	{Timever{0, 0, 1, ""}, Timever{0, 0, 1, ""}, 0},
	{Timever{0, 0, 2, ""}, Timever{0, 0, 1, ""}, 1},
	{Timever{1, 2, 3, ""}, Timever{1, 2, 3, ""}, 0},
	{Timever{2, 2, 4, ""}, Timever{1, 2, 4, ""}, 1},
	{Timever{1, 3, 3, ""}, Timever{1, 2, 3, ""}, 1},
	{Timever{1, 2, 4, ""}, Timever{1, 2, 3, ""}, 1},
	{Timever{1, 0, 0, ""}, Timever{2, 0, 0, ""}, -1},
	{Timever{2, 0, 0, ""}, Timever{2, 1, 0, ""}, -1},
	{Timever{2, 1, 0, ""}, Timever{2, 1, 1, ""}, -1},
	{Timever{1, 0, 0, "ce"}, Timever{1, 0, 0, "ee"}, 0},
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

type wrongTimeverFormatTest struct {
	str string
}

var wrongTimeverFormatTests = []wrongTimeverFormatTest{
	{""},
	{"."},
	{"1."},
	{".1"},
	{"a.b.c"},
	{"1.a.b"},
	{"1.1.a"},
	{"1.a.1"},
	{"a.1.1"},
	{".."},
	{"1.."},
	{"1.1."},
	{"1..1"},
	{"-1.1.1"},
	{"1.-1.1"},
	{"1.1.-1"},
	{"100.1.1"},
	{"1.13.1"},
}

func TestTimeverWrongFormat(t *testing.T) {
	for _, test := range wrongTimeverFormatTests {
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
