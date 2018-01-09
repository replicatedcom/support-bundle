package graphql

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSupportBundleSpec(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))
	defer ts.Close()

	sbs := &SupportBundleSpec{
		CustomerID: "123456",
		Endpoint:   ts.URL,
	}

	specs, err := sbs.Get()

	if err != nil {
		t.Fail()
	}
}
