package retraced

import (
	"fmt"
	"os"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

var _ = PDescribe("The retraced.events spec", func() {
	opts := types.RetracedEventsOptions{
		RetracedAPIClientOptions: types.RetracedAPIClientOptions{
			APIEndpoint: os.Getenv("RETRACED_API_ENDPOINT"),
			APIToken:    os.Getenv("RETRACED_API_TOKEN"),
			ProjectID:   os.Getenv("RETRACED_PROJECT_ID"),
			Insecure:    os.Getenv("RETRACED_INSECURE_SKIP_VERIFY") != "",
		},
	}

	BeforeEach(EnterNewTempDir)
	BeforeEach(func() {
		Expect(opts.APIEndpoint).NotTo(BeEmpty(), "RETRACED_API_ENDPOINT must be set")
		Expect(opts.APIToken).NotTo(BeEmpty(), "RETRACED_API_TOKEN must be set")
		Expect(opts.ProjectID).NotTo(BeEmpty(), "RETRACED_PROJECT_ID must be set")
	})
	AfterEach(CleanupDir)

	It("Collects events from retraced", func() {
		WriteBundleConfig(`
specs:
  - builtin: retraced.events
    raw: /audit/events.txt
    retraced.events:
        api_endpoint: ` + opts.APIEndpoint + `
        api_token: ` + opts.APIToken + `
        project_id: ` + opts.ProjectID + `
        insecure: ` + fmt.Sprintf("%v", opts.Insecure))

		GenerateBundle("--retraced")

		errors := GetFileFromBundle("error.json")
		Expect(errors).To(Equal("null"))

		contents := GetFileFromBundle("audit/events.txt")
		Expect(contents).ToNot(BeEmpty())

		header := strings.Split(contents, "\n")[0]
		Expect(header).To(Equal("action,crud,canonical_time"))

	})

	It("Allows for a custom mask+query", func() {
		WriteBundleConfig(`
specs:
  - builtin: retraced.events
    raw: /audit/events.txt
    retraced.events:
        api_endpoint: ` + opts.APIEndpoint + `
        api_token: ` + opts.APIToken + `
        project_id: ` + opts.ProjectID + `
        insecure: ` + fmt.Sprintf("%v", opts.Insecure) + `
        mask:
          CRUD: true
        query:
          CRUD: r`)

		GenerateBundle()

		errors := GetFileFromBundle("error.json")
		Expect(errors).To(Equal("null"))

		contents := GetFileFromBundle("audit/events.txt")
		Expect(contents).ToNot(BeEmpty())

		lines := strings.Split(contents, "\n")
		Expect(lines[0]).To(Equal("crud"))
		for _, line := range lines[1:] {
			if line != "" {
				Expect(line).To(Equal("r"))
			}
		}

	})
})
