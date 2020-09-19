package vulcain

// Adapted from the Hades project (https://github.com/gabesullice/hades/blob/master/lib/server/pusher.go)
// Copyright (c) 2019 Gabriel Sullice
// MIT License

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

const internalRequestHeader = "Vulcain-Explicit-Request"

// From the RFC:
//   The server SHOULD send PUSH_PROMISE (Section 6.6) frames prior to sending any frames that reference the promised responses.
//   This avoids a race where clients issue requests prior to receiving any PUSH_PROMISE frames.
type waitPusher struct {
	id         string
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

func newWaitPusher(p http.Pusher, id string, maxPushes int) *waitPusher {
	return &waitPusher{
		internalPusher: p,
		id:             id,
		maxPushes:      maxPushes,
		pushedURLs:     make(map[string]struct{}),
	}
}

type pushers struct {
	sync.RWMutex
	maxPushes int
	pusherMap map[string]*waitPusher
	logger    *zap.Logger
}

func (p *pushers) add(w *waitPusher) {
	p.Lock()
	defer p.Unlock()
	p.pusherMap[w.id] = w
}

func (p *pushers) get(id string) *waitPusher {
	p.RLock()
	defer p.RUnlock()
	return p.pusherMap[id]
}

func (p *pushers) remove(id string) {
	p.Lock()
	defer p.Unlock()
	delete(p.pusherMap, id)
}

// End of the code adapted from the Hades project

// Copyright (c) 2020 Kévin Dunglas
// APGLv3 License

func (p *pushers) getPusherForRequest(rw http.ResponseWriter, req *http.Request) (w *waitPusher) {
	internalPusher, ok := rw.(http.Pusher)
	if !ok {
		// Not an HTTP/2 connection
		return nil
	}

	// Need https://github.com/golang/go/issues/20566 to get rid of this hack
	explicitRequestID := req.Header.Get(internalRequestHeader)
	if explicitRequestID != "" {
		w = p.get(explicitRequestID)
		if w == nil {
			// Should not happen, is an attacker forging an evil request?
			p.logger.Debug("pusher not found", zap.String("url", req.RequestURI), zap.String("explicitRequestID", explicitRequestID))
			req.Header.Del(internalRequestHeader)
			explicitRequestID = ""
		}
	}

	if explicitRequestID == "" {
		// Explicit request
		w = newWaitPusher(internalPusher, uuid.Must(uuid.NewV4()).String(), p.maxPushes)
		p.add(w)
	}

	return w
}

func (p *pushers) cleanupAfterRequest(req *http.Request, w *waitPusher) {
	if w == nil {
		return
	}

	if req.Header.Get(internalRequestHeader) != "" {
		w.Done()
		return
	}

	// Wait for subrequests to finish
	w.Wait()
	p.remove(w.id)
}
