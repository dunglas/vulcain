package vulcain

// Adapted from the Hades project (https://github.com/gabesullice/hades/blob/master/lib/server/pusher.go)
// Copyright (c) 2019 Gabriel Sullice
// MIT License

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

const internalRequestHeader = "Vulcain-Explicit-Request"

type ctxKey struct{}

// waitPusher pushes relations and allow to wait for all PUSH_PROMISE to be sent
// From the RFC:
//   The server SHOULD send PUSH_PROMISE (Section 6.6) frames prior to sending any frames that reference the promised responses.
//   This avoids a race where clients issue requests prior to receiving any PUSH_PROMISE frames.
//
// Use newWaitPusher() to create a wait pusher
type waitPusher struct {
	id         string
	nbPushes   int
	pushedURLs map[string]struct{}
	maxPushes  int
	sync.WaitGroup
	sync.RWMutex
	internalPusher http.Pusher
}

// errRelationAlreadyPushed occurs when the relation has already been pushed
var errRelationAlreadyPushed = errors.New("relation already pushed")

func (p *waitPusher) Push(url string, opts *http.PushOptions) error {
	if p.maxPushes != -1 && p.nbPushes >= p.maxPushes {
		return fmt.Errorf("maximum allowed pushes (%d) reached", p.maxPushes)
	}

	cacheKey := fmt.Sprintf(":p:%v:f:%v:u:%s", opts.Header["Preload"], opts.Header["Fields"], url)

	p.Lock()
	if _, ok := p.pushedURLs[cacheKey]; ok {
		p.Unlock()
		return errRelationAlreadyPushed
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

// newWaitPusher creates a new waitPusher
func newWaitPusher(p http.Pusher, id string, maxPushes int) *waitPusher {
	return &waitPusher{
		internalPusher: p,
		id:             id,
		maxPushes:      maxPushes,
		pushedURLs:     make(map[string]struct{}),
	}
}

// pushers stores the list of current active pusher
// The same pusher is shared for the explicit response and all pushed responses
type pushers struct {
	sync.RWMutex
	maxPushes int
	pusherMap map[string]*waitPusher
	logger    *zap.Logger
}

// add adds a new waitPusher to the list
func (p *pushers) add(w *waitPusher) {
	p.Lock()
	defer p.Unlock()
	p.pusherMap[w.id] = w
}

// get gets the waitPusher from the list
func (p *pushers) get(id string) *waitPusher {
	p.RLock()
	defer p.RUnlock()
	return p.pusherMap[id]
}

// remove removes the waitPusher from the list
func (p *pushers) remove(id string) {
	p.Lock()
	defer p.Unlock()
	delete(p.pusherMap, id)
}

// End of the code adapted from the Hades project

// Copyright (c) 2020 Kévin Dunglas
// APGLv3 License

// getPusherForRequest retrieves the pusher associated with the explicit request
func (p *pushers) getPusherForRequest(rw http.ResponseWriter, req *http.Request) (w *waitPusher) {
	internalPusher, ok := rw.(http.Pusher)
	if !ok {
		// Not an HTTP/2 connection
		return nil
	}

	// Need https://github.com/golang/go/issues/20566 to get rid of this hack
	explicitRequestID := req.Header.Get(internalRequestHeader)
	if explicitRequestID == "" {
		// This is the explicit request, let's create a wait pusher
		w = newWaitPusher(internalPusher, uuid.Must(uuid.NewV4()).String(), p.maxPushes)
		p.add(w)

		return w
	}

	if w = p.get(explicitRequestID); w != nil {
		return w
	}

	// Should not happen, is an attacker forging an evil request?
	p.logger.Debug("pusher not found", zap.String("url", req.RequestURI), zap.String("explicitRequestID", explicitRequestID))
	req.Header.Del(internalRequestHeader)

	return nil
}

// finish waits for all PUSH_PROMISEs to be sent before returning for the explicit request.
func (p *pushers) finish(req *http.Request, wait bool) {
	pusher := req.Context().Value(ctxKey{}).(*waitPusher)
	if pusher == nil {
		return
	}

	if req.Header.Get(internalRequestHeader) != "" {
		pusher.Done()
		return
	}

	// Wait for subrequests to finish, except if it's an error to release resources as soon as possible
	if wait {
		pusher.Wait()
	}

	p.remove(pusher.id)
}
