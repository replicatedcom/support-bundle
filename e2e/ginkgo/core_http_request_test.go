package ginkgo

import (
	"fmt"
	"net"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("os.http-request", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

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
  - os.http-request:
      url: http://%s/raw
    output_dir: /os/http-request/raw/
  - os.http-request:
      url: http://%s/json
    output_dir: /os/http-request/json/
  - os.http-request:
      url: http://%s/post
      method: post
    output_dir: /os/http-request/post/
  - os.http-request:
      url: http://%s/post
    output_dir: /os/http-request/err/`,
				ln.Addr(), ln.Addr(), ln.Addr(), ln.Addr()))

			GenerateBundle()

			_ = GetResultFromBundle("os/http-request/raw/body")
			contents := GetFileFromBundle("os/http-request/raw/body")
			Expect(contents).To(Equal("Hello World!"))

			_ = GetResultFromBundle("os/http-request/json/body")
			contents = GetFileFromBundle("os/http-request/json/body")
			Expect(contents).To(Equal(`{"a":"b"}`))

			_ = GetResultFromBundle("os/http-request/post/body")
			contents = GetFileFromBundle("os/http-request/post/body")
			Expect(contents).To(Equal(`aok!`))

			ExpectBundleErrorToHaveOccured("os/http-request/err", "unexpected status 405 Method Not Allowed")
		})
	})

	/*	It("Successfully executes the https request", func() {

				x509Cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
				Expect(err).NotTo(HaveOccurred())
				ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
					Certificates: []tls.Certificate{x509Cert},
				})
				Expect(err).NotTo(HaveOccurred())
				defer ln.Close()
				mux := http.NewServeMux()
				mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`Hello World!`))
				})
				go http.Serve(ln, mux)

				WriteBundleConfig(fmt.Sprintf(`
		specs:
		  - builtin: core.http-request
		    raw: /core/http-request/insecure
		    core.http-request:
		      url: https://%s/raw
		      insecure: true
		  - builtin: core.http-request
		    raw: /core/http-request/secure
		    core.http-request:
		      url: https://%s/raw`,
					ln.Addr(), ln.Addr()))

				GenerateBundle()

				contents := GetFileFromBundle("core/http-request/insecure")
				Expect(strings.TrimSpace(contents)).To(Equal("Hello World!"))

				ExpectBundleErrorToHaveOccured(
					"core/http-request/secure",
					`make request: Get https:\/\/.+\/raw: x509: cannot validate certificate .+`)
			})
	*/
})

var certPEM, keyPEM = `-----BEGIN CERTIFICATE-----
MIIDqTCCApGgAwIBAgIJAKtAr29xIpOgMA0GCSqGSIb3DQEBCwUAMEExEzARBgNV
BAMTCmRvbWFpbi5jb20xHTAbBgNVBAoTFE15IENvbXBhbnkgTmFtZSBMVEQuMQsw
CQYDVQQGEwJVUzAeFw0xNzExMDcyMTQwNDlaFw0xODExMDcyMTQwNDlaMEExEzAR
BgNVBAMTCmRvbWFpbi5jb20xHTAbBgNVBAoTFE15IENvbXBhbnkgTmFtZSBMVEQu
MQswCQYDVQQGEwJVUzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOxw
p5h/0XOPAlRyCjd49NT0nMUM6ZIOo6udEBRIHuZEg6lCvZl2MhOgcdy1e6Uhbr8G
vrRtnyeQW1g8qwooZqSQyRheO8Q0+wkQduCl9sSZ1cX6eT7L1uv/VCj1x3hkJx+/
gB+w4cFuWvAoVDHCVFpVHoOnj+ApcDQRdErSwt8t2V+MvxK+trjrMCNH77i99dGL
UQIvrzt09MlRTdRdmzN2FqR64ApWJxtG4TR4K+N3oujfSy68QKo/XRrbcS3Fdbqm
zdlxTFWX4oByFGg3oKMTmOqF5ncmy0ntBO0O6+IbQYGXVpXWr0i1zg+nE35vYhHv
qIcFUd84QoVBPss0LTMCAwEAAaOBozCBoDAdBgNVHQ4EFgQUuIsB6y9CEoy6KVj8
q10dnHnog00wcQYDVR0jBGowaIAUuIsB6y9CEoy6KVj8q10dnHnog02hRaRDMEEx
EzARBgNVBAMTCmRvbWFpbi5jb20xHTAbBgNVBAoTFE15IENvbXBhbnkgTmFtZSBM
VEQuMQswCQYDVQQGEwJVU4IJAKtAr29xIpOgMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADggEBANnWCeDcEcM8jCU9XMPN2NThf6bUpQW0C3p017fOBtVrc6N2
Q6eGdYK7AtXRGfjpicmAAL19mEa9oaDnn9pk2k1F5RZeWs4/ztUMknkJpUqqmYsh
ziiN6Ze6eXc1v1TvdKxQh9latDeUZVASRbFlQRlTzGmzhkiIPPNi/1M+VVp+IrVW
5tW4TCpEsUWzUA09CylSWa2OvhXE8FCfVZIroieoRredU+yIG9Q/OzSEhqX2zPDI
jp8tR224aCy2YCfpBNK1gYpRnx5H1D726plifzDEYPl6dEx9DkttglPib7BC1kwm
NcSFYOm5gwgZ+qmP3AgeYwB8+FaEJoRcfGF2ltI=
-----END CERTIFICATE-----`, `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA7HCnmH/Rc48CVHIKN3j01PScxQzpkg6jq50QFEge5kSDqUK9
mXYyE6Bx3LV7pSFuvwa+tG2fJ5BbWDyrCihmpJDJGF47xDT7CRB24KX2xJnVxfp5
PsvW6/9UKPXHeGQnH7+AH7DhwW5a8ChUMcJUWlUeg6eP4ClwNBF0StLC3y3ZX4y/
Er62uOswI0fvuL310YtRAi+vO3T0yVFN1F2bM3YWpHrgClYnG0bhNHgr43ei6N9L
LrxAqj9dGttxLcV1uqbN2XFMVZfigHIUaDegoxOY6oXmdybLSe0E7Q7r4htBgZdW
ldavSLXOD6cTfm9iEe+ohwVR3zhChUE+yzQtMwIDAQABAoIBAA73AmSQkn0x6//U
R/lC2pwv70w5iP8JlipigoYGGo6Qf5TS+JUh+gFsAkjp462L9Cp7Ds70sUIbzOxE
yr6V2AlKcK/uJvJAODNWq2+EkZ3X7sPdGpdy52OAgZ8mDz37eM51nHJlC6TmizLw
GoQbvKMLCCdlD6IsfUGOaUHjHRpq+VUpcLZ4NYh/GzBOyQHHJo24jFIoZbr5hvJG
dPj5uMXFJjN+pJUzgeBO9Zj93zans6yWQokgYOHX3arCFxOjniH/ype34ymen1sR
2qRcvCy/7gTIv9Hxu5EknRXCp5Wt/Ig1Agi9vYhsHAjNTl7Zjhc4bNUVM4n/EKTW
iyNyhiECgYEA+l2fdHEWMduNK9eTknZ/kpwAD1LZBf02AcV3Qm1I9l8qm8ymxcEs
8woxMrfEL7BMHmnAHfFzSh1ot4ZG6NxetJlzjW0lXAXnU1gTx7bRIUxbo5xzLDsD
9NVAWNJc3cQnxOeSO2QW6zO+2hE2E/12Efe/d1YgICa8on7eJqqQAu0CgYEA8cLO
EhrQBHNsG4mLzBEbjxj2hGsowD2NFlQ+UYww7PIMCX3qzz+HEcOEaCIi0LgNAmFg
5a+mpe9KtDr7m/XlClld7EWh6smIA1O9+lmN9nu99YH9tOJti3pXXbOFiUhQqmIi
GRyIU98L6JTyaCpUGf2ktao6UxjYxwknG3X0TJ8CgYBQJ+VDLGmEsNvzq2Mtww54
68UBIu8kgbrmukfCVqbDahiEJPNH4N75OMwjhr4i3nigTA8cBw94LQ43o5/UMamI
fJCIOOd7HNDA2DQM/rTZyk6UhSRChupvWk7toPvmbESnP9SLezHzP2/c9SGxKLbC
beU42bQTVxORmriY/IZ6yQKBgQCcfsSUNaUH7ItDfBLxYvWa+MbCyvcTEgTdOmUo
tn4JM1mVX1v7Eh1l41E3czlkMG/DZbOqmrxeV3rdFf0/ZLoBq/2/bwe0CwavWKr2
frgFoO5DGQVY7OWKTwR01DuRtSz6ThHSfYTF/fEgeiI8SYItXOIc8ndUyRWyKXuW
LBGa8wKBgQCKXvjud0839HE4dn8FWknSKY8CrTBj/Tl7Inb5FqQVFz0wKV4EeVGm
kGFcpGiX3r6ncmfjP0+CAbqSU1Ks2GB64+lRXCwSKlkgkMkf92vesylQFVH73qFn
uu2S3QvhJ/1n87APPtbVFFJ3BsfSKyH95GJoNY/IPtN8ZVMlb/zGBw==
-----END RSA PRIVATE KEY-----`
