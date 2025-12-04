package comm

import "sync"

type Single[T any] struct {
	once sync.Once
	obj  *T
}

func (s *Single[T]) Instance(newFun ...func() *T) *T {
	s.once.Do(func() {
		if len(newFun) != 0 {
			s.obj = newFun[0]()
		} else {
			s.obj = new(T)
		}
	})

	return s.obj
}
