package gateway

// Adapted from the Hades project (https://github.com/gabesullice/hades/blob/master/lib/server/pusher.go)
// Copyright (c) 2019 Gabriel Sullice
// MIT License

import (
	"fmt"
	"net/http"
	"sync"
)

// From the RFC:
//   The server SHOULD send PUSH_PROMISE (Section 6.6) frames prior to sending any frames that reference the promised responses.
//   This avoids a race where clients issue requests prior to receiving any PUSH_PROMISE frames.
type waitPusher struct {
	nbPushes   int
	pushedURLs map[string]struct{}
	maxPushes  int
	sync.WaitGroup
	sync.RWMutex
	internalPusher http.Pusher
}

type relationAlreadyPushedError struct{}

func (f relationAlreadyPushedError) Error() string {
	return "Relation already pushed"
}

func (p *waitPusher) Push(url string, opts *http.PushOptions) error {
	if p.maxPushes != -1 && p.nbPushes >= p.maxPushes {
		return fmt.Errorf("Maximum allowed pushes (%d) reached", p.maxPushes)
	}

	cacheKey := fmt.Sprintf(":p:%v:f:%v:u:%s", opts.Header["Preload"], opts.Header["Fields"], url)

	p.Lock()
	if _, ok := p.pushedURLs[cacheKey]; ok {
		p.Unlock()
		return &relationAlreadyPushedError{}
	}

	p.nbPushes++
	p.pushedURLs[cacheKey] = struct{}{}
	p.Unlock()

	p.Add(1)
	if err := p.internalPusher.Push(url, opts); err != nil {
		p.Done()
		return err
	}

	return nil
}

func newWaitPusher(p http.Pusher, maxPushes int) *waitPusher {
	return &waitPusher{
		internalPusher: p,
		maxPushes:      maxPushes,
		pushedURLs:     make(map[string]struct{}),
	}
}

type pushers struct {
	sync.RWMutex
	pusherMap map[string]*waitPusher
}

func (p *pushers) add(id string, w *waitPusher) {
	p.Lock()
	p.pusherMap[id] = w
	p.Unlock()
}

func (p *pushers) get(id string) (*waitPusher, bool) {
	p.RLock()
	w, ok := p.pusherMap[id]
	p.RUnlock()
	return w, ok
}

func (p *pushers) remove(id string) {
	p.Lock()
	delete(p.pusherMap, id)
	p.Unlock()
}
