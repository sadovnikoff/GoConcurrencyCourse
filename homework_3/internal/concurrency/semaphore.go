package concurrency

type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		tickets: make(chan struct{}, limit),
	}
}

func (s *Semaphore) Acquire() {
	if s == nil || s.tickets == nil {
		return
	}

	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	if s == nil || s.tickets == nil {
		return
	}

	<-s.tickets
}
