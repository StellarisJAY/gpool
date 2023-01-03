package gpool

type workerQueue interface {
	poll() *worker
	put(w *worker)
	len() int
	cap() int
}

func newWorkerQueue(queueType WorkerQueueType, cap int32) workerQueue {
	switch queueType {
	case SliceWorkerQueue:
		return newSliceQueue(cap)
	case RingWorkerQueue:
	case LIFOWorkerQueue:
	default:

	}
	return nil
}
