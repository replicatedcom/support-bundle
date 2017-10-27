package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("Scrubbing secrets from file", func() {
	BeforeEach(func() {
		log.Println("lol")
	})
	Specify("Then Replicated should redirect to the new hostname", func() {
		Expect("foo").Should(Equal("foo"))

	})
});
