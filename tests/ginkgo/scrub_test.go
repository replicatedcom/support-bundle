package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scrubbing Secrets", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Scrubs any instances of PGPASSWORD=.* when reading local file (ByteSource)", func() {

		WriteFile("pg.env", `
PGDATABASE=mydata
PGPASSWORD=mypass`)

		WriteBundleConfig(`
specs:
  - builtin: core.read-file
    raw: /pg/pg.env
    config:
      file_path: pg.env
      scrub:
        regex: (PGPASSWORD)=(.*)
        replace: $1=REDACTED
      `)

		GenerateBundle()

		contents := GetFileFromBundle("pg/pg.env")

		Expect(contents).To(ContainSubstring("PGDATABASE=mydata"))
		Expect(contents).NotTo(ContainSubstring("PGPASSWORD=mypass"))
		Expect(contents).To(ContainSubstring("PGPASSWORD=REDACTED"))
	})

	Describe("Container tests", func() {
		var containerID string
		BeforeEach(func() {
			containerID = MakeDockerContainer()
		})
		AfterEach(func() {
			RemoveDockerContainer(containerID)
		})

		It("Scrubs any instances of [uU]buntu from file on docker container (StreamSource)", func() {
			WriteBundleConfig(`
specs:
- builtin: docker.read-file
  raw: /os/release-uncensored
  config:
    file_path: /usr/lib/os-release
    container_id: ` + containerID + `
- builtin: docker.read-file
  raw: /os/release
  config:
    file_path: /usr/lib/os-release
    scrub:
      regex: "[uU]buntu"
      replace: replicatedOS
    container_id: ` + containerID)

			GenerateBundle()

			censored := GetFileFromBundle("os/release")
			uncensored := GetFileFromBundle("os/release-uncensored")

			Expect(censored).To(ContainSubstring("replicatedOS"))
			Expect(censored).NotTo(ContainSubstring("ubuntu"))

			Expect(uncensored).NotTo(ContainSubstring("replicatedOS"))
			Expect(uncensored).To(ContainSubstring("ubuntu"))
		})

		It("Scrubs any instances of [uU]buntu from command output run on docker container (StreamsSource)", func() {
			WriteBundleConfig(`
specs:
- builtin: docker.exec-command
  raw: /os/release-uncensored.
  config:
    command: cat
    args: ["/usr/lib/os-release"]
    container_id: ` + containerID + `
- builtin: docker.exec-command
  raw: /os/release.
  config:
    command: cat
    args: ["/usr/lib/os-release"]
    scrub:
      regex: "[uU]buntu"
      replace: replicatedOS
    container_id: ` + containerID)

			GenerateBundle()

			censored := GetFileFromBundle("os/release.stdout")
			uncensored := GetFileFromBundle("os/release-uncensored.stdout")

			Expect(censored).To(ContainSubstring("replicatedOS"))
			Expect(censored).NotTo(ContainSubstring("ubuntu"))

			Expect(uncensored).NotTo(ContainSubstring("replicatedOS"))
			Expect(uncensored).To(ContainSubstring("ubuntu"))
		})
	})
})
