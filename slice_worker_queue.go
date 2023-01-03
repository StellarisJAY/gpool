package gpool

type sliceQueue struct {
	cache    []*worker
	capacity int32
}

func newSliceQueue(cap int32) *sliceQueue {
	return &sliceQueue{
		cache:    make([]*worker, 0, cap),
		capacity: cap,
	}
}

func (s *sliceQueue) poll() *worker {
	if s.len() == 0 {
		return nil
	}
	w := s.cache[0]
	s.cache = s.cache[1:]
	return w
}

func (s *sliceQueue) put(w *worker) {
	if s.len() < s.cap() {
		s.cache = append(s.cache, w)
	}
}

func (s *sliceQueue) len() int {
	return len(s.cache)
}

func (s *sliceQueue) cap() int {
	return int(s.capacity)
}
