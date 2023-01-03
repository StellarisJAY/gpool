package gpool

type WorkerQueueType byte

const (
	minTaskCapacity = 1

	SliceWorkerQueue WorkerQueueType = iota
	RingWorkerQueue
	LIFOWorkerQueue
)

type Options struct {
	taskCapacity int32
	poolCapacity int32
	panicHandler func(p any)

	queueType WorkerQueueType
}
