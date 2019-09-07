package gateway

// Adapted from https://github.com/gabesullice/hades/blob/master/lib/server/pusher.go
// Copyright (c) 2019 Gabriel Sullice
// MIT License

import (
	"net/http"
	"sync"
)

// From the RFC:
//   The server SHOULD send PUSH_PROMISE (Section 6.6) frames prior to sending any frames that reference the promised responses.
//   This avoids a race where clients issue requests prior to receiving any PUSH_PROMISE frames.
type pusher struct {
	sync.WaitGroup
	internalPusher http.Pusher
}

func (r *pusher) Push(target string, opts *http.PushOptions) error {
	r.Add(1)
	if err := r.internalPusher.Push(target, opts); err != nil {
		r.Done()

		return err
	}

	return nil
}

type pushers struct {
	sync.RWMutex
	pusherMap map[string]*pusher
}

func (r *pushers) add(id string, p *pusher) {
	r.Lock()
	r.pusherMap[id] = p
	r.Unlock()
}

func (r *pushers) get(id string) (*pusher, bool) {
	r.RLock()
	p, ok := r.pusherMap[id]
	r.RUnlock()
	return p, ok
}

func (r *pushers) remove(id string) {
	r.Lock()
	delete(r.pusherMap, id)
	r.Unlock()
}
