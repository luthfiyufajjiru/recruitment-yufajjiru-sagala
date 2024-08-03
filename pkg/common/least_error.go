package common

import "sync"

type (
	LeastError struct {
		mtx sync.RWMutex
		err error
	}
)

func (e *LeastError) setError(err error) {
	defer e.mtx.Unlock()
	e.mtx.Lock()
	e.err = err
}

func (e *LeastError) Do(inp func() (err error)) {
	defer e.mtx.RUnlock()
	e.mtx.RLock()
	if err := inp(); e.err == nil && err != nil {
		e.mtx.RUnlock()
		e.setError(err)
		e.mtx.RLock()
	}
}

func (e *LeastError) Err() error {
	defer e.mtx.RUnlock()
	e.mtx.RLock()
	return e.err
}
