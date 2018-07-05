package producers

import (
	"context"
	"io"

	"path/filepath"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/exec"

	jww "github.com/spf13/jwalterweatherman"
)

// ContainerCp copies a file/folder from the provided directory within the pod+container.
// if container=="", it copies that file/folder from all containers in the pod.
func (k *Kubernetes) ContainerCp(pod, container, namespace, path string) types.StreamsProducer {
	ns := namespace
	if ns == "" {
		ns = metav1.NamespaceDefault
	}

	containerNameProvided := container != ""

	return func(ctx context.Context) (map[string]io.Reader, error) {
		//if container == "", get all containers in pod & loop over them

		var containers []string

		if containerNameProvided {
			containers = []string{container}
		} else {
			thisPod, err := k.client.CoreV1().Pods(ns).Get(pod, metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "get containers in pod")
			}
			for _, containerStatus := range thisPod.Status.ContainerStatuses {
				if containerStatus.State.Running != nil {
					containers = append(containers, containerStatus.Name)
				}
			}
		}

		retMap := make(map[string]io.Reader)

		for _, cont := range containers {
			reader, outStream := io.Pipe()

			req := k.client.CoreV1().RESTClient().Post().
				Resource("pods").Name(pod).Namespace(ns).
				SubResource("exec").Param("container", cont)

			req.VersionedParams(&corev1.PodExecOptions{
				Container: cont,
				Command:   []string{"tar", "cf", "-", "-C", filepath.Dir(path), filepath.Base(path)},
				Stdin:     false,
				Stdout:    true,
				Stderr:    false,
				TTY:       false,
			}, scheme.ParameterCodec)

			executor, err := remotecommand.NewSPDYExecutor(k.config, "POST", req.URL())

			if err != nil {
				return nil, errors.Wrap(err, "NewExecutor")
			}

			go func(outStream io.WriteCloser) {
				//TODO: error handling
				if err := executor.Stream(remotecommand.StreamOptions{
					Stdin:  nil,
					Stdout: outStream,
					Stderr: nil,
					Tty:    false,
				}); err != nil {
					if _, ok := err.(exec.CodeExitError); ok {
						//handle exit code error (forward stderr?)
					}
					jww.ERROR.Printf("Got error: %s", err.Error())
				}
			}(outStream)

			retMap[cont] = reader
		}

		return retMap, nil
	}
}
