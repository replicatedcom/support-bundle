package producers

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_logsReaderWithTimeout(t *testing.T) {
	pr, pw := io.Pipe()
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			pw.Write([]byte(fmt.Sprintf("%d\n", i)))
		}
		pw.Close()
	}()
	ppr := logsReaderWithTimeout(pr, time.Second)
	b, err := ioutil.ReadAll(ppr)
	assert.NoError(t, err)
	assert.Equal(t, "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n", string(b))
}

func Test_logsReaderWithTimeoutTimeout(t *testing.T) {
	pr, pw := io.Pipe()
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(11*i) * time.Millisecond)
			pw.Write([]byte(fmt.Sprintf("%d\n", i)))
		}
		pw.Close()
	}()
	ppr := logsReaderWithTimeout(pr, 50*time.Millisecond)
	start := time.Now()
	b, err := ioutil.ReadAll(ppr)
	end := time.Now().Sub(start)
	assert.Equal(t, fmt.Errorf("reader timeout after %s", 50*time.Millisecond), err)
	assert.Equal(t, "0\n1\n2\n3\n4\n", string(b))
	assert.Condition(t, func() bool {
		return end > 160*time.Millisecond
	}, "read took less than %s (%s)", 160*time.Millisecond, end)
}
