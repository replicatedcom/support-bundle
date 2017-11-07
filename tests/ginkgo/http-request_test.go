package ginkgo

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Make http request", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	It("Successfully executes the http request", func() {

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		defer ln.Close()
		mux := http.NewServeMux()
		mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`Hello World!`))
		})
		mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"a":"b"}`))
		})
		mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aok!`))
		})
		go http.Serve(ln, mux)

		WriteBundleConfig(fmt.Sprintf(`
specs:
  - builtin: core.http-request
    raw: /core/http-request/raw
    core.http-request:
      URL: http://%s/raw
  - builtin: core.http-request
    json: /core/http-request/json
    core.http-request:
      URL: http://%s/json
  - builtin: core.http-request
    raw: /core/http-request/post
    core.http-request:
      URL: http://%s/post
      Method: post
  - builtin: core.http-request
    raw: /core/http-request/err
    core.http-request:
      URL: http://%s/post`,
			ln.Addr(), ln.Addr(), ln.Addr(), ln.Addr()))

		GenerateBundle()

		contents := GetFileFromBundle("core/http-request/raw")
		Expect(strings.TrimSpace(contents)).To(Equal("Hello World!"))

		contents = GetFileFromBundle("core/http-request/json")
		Expect(strings.TrimSpace(contents)).To(Equal(`{"a":"b"}`))

		contents = GetFileFromBundle("core/http-request/post")
		Expect(strings.TrimSpace(contents)).To(Equal(`aok!`))

		ExpectBundleErrorToHaveOccured("core/http-request/err", "unexpected status 405 Method Not Allowed")
	})

})
