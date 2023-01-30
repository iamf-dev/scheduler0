package workers

import (
	"testing"
)

func Test_Workers(t *testing.T) {

	t.Run("should execute callback when work is queued", func(t *testing.T) {
		pool := make(chan chan []interface{}, 1)

		callback := func(effector func(s chan any, e chan any), s chan any, e chan any) {
			effector(s, e)
		}

		worker := NewWorker(pool, callback)

		worker.Start()

		e := make(chan any, 1)
		s := make(chan any, 1)
		c := func(s chan any, e chan any) {
			close(s)
			close(e)
		}

		worker.WorkerQueue <- []interface{}{c, s, e}

		_, eClosed := <-e
		_, sClosed := <-s

		if eClosed || sClosed {
			t.Fatal("should execute callback and close channel")
		}
	})

}
