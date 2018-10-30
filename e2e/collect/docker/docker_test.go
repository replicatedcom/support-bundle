package docker

import (
	"context"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
	uuid "github.com/satori/go.uuid"

	"github.com/replicatedcom/support-bundle/pkg/collect/plugins/docker/planners"
)

var _ = Describe("docker container select", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	var containerID1, containerName1 string
	var containerID2, containerName2 string
	var containerID3, containerName3 string
	containerName1 = uuid.NewV4().String()
	containerName2 = uuid.NewV4().String()
	containerName3 = uuid.NewV4().String()
	BeforeEach(func() {
		containerID1 = MakeDockerContainer(dockerClient, containerName1, map[string]string{"label1": "one", "label2": "both"}, nil)
		containerID2 = MakeDockerContainer(dockerClient, containerName2, map[string]string{"label1": "two", "label2": "both"}, nil)
		containerID3 = MakeDockerContainer(dockerClient, containerName3, nil, nil)
	})
	AfterEach(func() {
		RemoveDockerContainer(dockerClient, containerID1)
		RemoveDockerContainer(dockerClient, containerID2)
		RemoveDockerContainer(dockerClient, containerID3)
	})

	Context("When selecting containers by name or labels", func() {

		It("should return the correct container ID", func() {

			var containers []string
			var err error

			d := planners.New(nil, dockerClient)

			containers, err = d.SelectContainerIDs("", []string{"label1=one"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(containers)).To(Equal(1))
			Expect(containers).To(ContainElement(containerID1))

			containers, err = d.SelectContainerIDs("", []string{"label2=both"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(containers)).To(Equal(2))
			Expect(containers).To(ContainElement(containerID1))
			Expect(containers).To(ContainElement(containerID2))

			containers, err = d.SelectContainerIDs(containerName3, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(containers)).To(Equal(1))
			Expect(containers).To(ContainElement(containerName3))
		})
	})
})
