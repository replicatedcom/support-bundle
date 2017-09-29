package bundle

import (
	"time"

	"github.com/replicatedcom/support-bundle/types"
)

type Task struct {
	Description string
	ExecFunc    func(chan types.Data, chan bool, chan types.Result, time.Duration, []string) error
	Timeout     time.Duration
	Args        []string
}

// type Data struct {
// 	Filename string
// 	Data     interface{}
// }
