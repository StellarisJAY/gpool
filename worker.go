package gpool

type worker struct {
	pool  *Pool
	tasks chan func()
}

func (w *worker) run() {
	w.pool.addRunning(1)
	go func() {
		defer func() {
			// task panic, call panicHandler
			if err := recover(); err != nil {
				if handler := w.pool.options.panicHandler; handler != nil {
					handler(err)
				} else {
					handlePanic(err)
				}
			}
			// todo return worker to worker queue
		}()
		for task := range w.tasks {
			// empty task stops worker
			if task == nil {
				return
			}
			task()
			w.pool.returnWorker(w)
		}
	}()
}

func handlePanic(p any) {
	switch p.(type) {
	case string:
	case error:
	default:

	}
}
