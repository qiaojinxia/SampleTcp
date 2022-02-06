package main

import (
	"fmt"
	"sync"
)

type HandleFunc func(sess *Session) error

type Router struct {
	lk     sync.RWMutex
	router map[uint32]HandleFunc
}

func (r *Router) RegisterRouter(no uint32, fn HandleFunc) error {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.router[no] = fn
	return nil
}

func (r *Router) GetHandler(no uint32) (HandleFunc, error) {
	r.lk.RLock()
	defer r.lk.RUnlock()
	if fn, exists := r.router[no]; !exists {
		return nil, fmt.Errorf("handler no:%d not exists", no)
	} else {
		return fn, nil
	}
}
