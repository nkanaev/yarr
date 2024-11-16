package nostr

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/puzpuzpuz/xsync/v3"
)

type Status int

var subscriptionIDCounter atomic.Int64

type Relay struct {
	closeMutex sync.Mutex

	URL           string
	RequestHeader http.Header // e.g. for origin header

	Connection    *Connection
	Subscriptions *xsync.MapOf[int64, *Subscription]

	ConnectionError         error
	connectionContext       context.Context // will be canceled when the connection closes
	connectionContextCancel context.CancelFunc

	challenge                     string       // NIP-42 challenge, we only keep the last
	noticeHandler                 func(string) // NIP-01 NOTICEs
	customHandler                 func([]byte) // nonstandard unparseable messages
	okCallbacks                   *xsync.MapOf[string, func(bool, string)]
	writeQueue                    chan writeRequest
	subscriptionChannelCloseQueue chan *Subscription

	// custom things that aren't often used
	//
	AssumeValid bool // this will skip verifying signatures for events received from this relay
}

type writeRequest struct {
	msg    []byte
	answer chan error
}

// NewRelay returns a new relay. The relay connection will be closed when the context is canceled.
func NewRelay(ctx context.Context, url string, opts ...RelayOption) *Relay {
	ctx, cancel := context.WithCancel(ctx)
	r := &Relay{
		URL:                           NormalizeURL(url),
		connectionContext:             ctx,
		connectionContextCancel:       cancel,
		Subscriptions:                 xsync.NewMapOf[int64, *Subscription](),
		okCallbacks:                   xsync.NewMapOf[string, func(bool, string)](),
		writeQueue:                    make(chan writeRequest),
		subscriptionChannelCloseQueue: make(chan *Subscription),
		RequestHeader:                 make(http.Header, 1),
	}

	for _, opt := range opts {
		opt.ApplyRelayOption(r)
	}

	return r
}

// RelayConnect returns a relay object connected to url.
// Once successfully connected, cancelling ctx has no effect.
// To close the connection, call r.Close().
func RelayConnect(ctx context.Context, url string, opts ...RelayOption) (*Relay, error) {
	r := NewRelay(context.Background(), url, opts...)
	err := r.Connect(ctx)
	return r, err
}

// When instantiating relay connections, some options may be passed.
// RelayOption is the type of the argument passed for that.
type RelayOption interface {
	ApplyRelayOption(*Relay)
}

var (
	_ RelayOption = (WithNoticeHandler)(nil)
	_ RelayOption = (WithCustomHandler)(nil)
)

// WithNoticeHandler just takes notices and is expected to do something with them.
// when not given, defaults to logging the notices.
type WithNoticeHandler func(notice string)

func (nh WithNoticeHandler) ApplyRelayOption(r *Relay) {
	r.noticeHandler = nh
}

// WithCustomHandler must be a function that handles any relay message that couldn't be
// parsed as a standard envelope.
type WithCustomHandler func(data []byte)

func (ch WithCustomHandler) ApplyRelayOption(r *Relay) {
	r.customHandler = ch
}

// String just returns the relay URL.
func (r *Relay) String() string {
	return r.URL
}

// Context retrieves the context that is associated with this relay connection.
func (r *Relay) Context() context.Context { return r.connectionContext }

// IsConnected returns true if the connection to this relay seems to be active.
func (r *Relay) IsConnected() bool { return r.connectionContext.Err() == nil }

// Connect tries to establish a websocket connection to r.URL.
// If the context expires before the connection is complete, an error is returned.
// Once successfully connected, context expiration has no effect: call r.Close
// to close the connection.
//
// The underlying relay connection will use a background context. If you want to
// pass a custom context to the underlying relay connection, use NewRelay() and
// then Relay.Connect().
func (r *Relay) Connect(ctx context.Context) error {
	return r.ConnectWithTLS(ctx, nil)
}

// ConnectWithTLS tries to establish a secured websocket connection to r.URL using customized tls.Config (CA's, etc).
func (r *Relay) ConnectWithTLS(ctx context.Context, tlsConfig *tls.Config) error {
	if r.connectionContext == nil || r.Subscriptions == nil {
		return fmt.Errorf("relay must be initialized with a call to NewRelay()")
	}

	if r.URL == "" {
		return fmt.Errorf("invalid relay URL '%s'", r.URL)
	}

	if _, ok := ctx.Deadline(); !ok {
		// if no timeout is set, force it to 7 seconds
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 7*time.Second)
		defer cancel()
	}

	if r.RequestHeader.Get("User-Agent") == "" {
		r.RequestHeader.Set("User-Agent", "github.com/nbd-wtf/go-nostr")
	}

	conn, err := NewConnection(ctx, r.URL, r.RequestHeader, tlsConfig)
	if err != nil {
		return fmt.Errorf("error opening websocket to '%s': %w", r.URL, err)
	}
	r.Connection = conn

	// ping every 29 seconds
	ticker := time.NewTicker(29 * time.Second)

	// to be used when the connection is closed
	go func() {
		<-r.connectionContext.Done()

		// stop the ticker
		ticker.Stop()

		// nil the connection
		r.Connection = nil

		// close all subscriptions
		for _, sub := range r.Subscriptions.Range {
			sub.Unsub()
		}
	}()

	// queue all write operations here so we don't do mutex spaghetti
	go func() {
		for {
			select {
			case <-ticker.C:
				err := wsutil.WriteClientMessage(r.Connection.conn, ws.OpPing, nil)
				if err != nil {
					InfoLogger.Printf("{%s} error writing ping: %v; closing websocket", r.URL, err)
					r.Close() // this should trigger a context cancelation
					return
				}
			case writeRequest := <-r.writeQueue:
				// all write requests will go through this to prevent races
				debugLogf("{%s} sending %v\n", r.URL, string(writeRequest.msg))
				if err := r.Connection.WriteMessage(r.connectionContext, writeRequest.msg); err != nil {
					writeRequest.answer <- err
				}
				close(writeRequest.answer)
			case <-r.connectionContext.Done():
				// stop here
				return
			}
		}
	}()

	// general message reader loop
	go func() {
		buf := new(bytes.Buffer)

		for {
			buf.Reset()
			if err := conn.ReadMessage(r.connectionContext, buf); err != nil {
				r.ConnectionError = err
				r.Close()
				break
			}

			message := buf.Bytes()
			debugLogf("{%s} %v\n", r.URL, message)
			envelope := ParseMessage(message)
			if envelope == nil {
				if r.customHandler != nil {
					r.customHandler(message)
				}
				continue
			}

			switch env := envelope.(type) {
			case *NoticeEnvelope:
				// see WithNoticeHandler
				if r.noticeHandler != nil {
					r.noticeHandler(string(*env))
				} else {
					log.Printf("NOTICE from %s: '%s'\n", r.URL, string(*env))
				}
			case *AuthEnvelope:
				if env.Challenge == nil {
					continue
				}
				r.challenge = *env.Challenge
			case *EventEnvelope:
				if env.SubscriptionID == nil {
					continue
				}

				if subscription, ok := r.Subscriptions.Load(subIdToSerial(*env.SubscriptionID)); !ok {
					// InfoLogger.Printf("{%s} no subscription with id '%s'\n", r.URL, *env.SubscriptionID)
					continue
				} else {
					// check if the event matches the desired filter, ignore otherwise
					if !subscription.match(&env.Event) {
						InfoLogger.Printf("{%s} filter does not match: %v ~ %v\n", r.URL, subscription.Filters, env.Event)
						continue
					}

					// check signature, ignore invalid, except from trusted (AssumeValid) relays
					if !r.AssumeValid {
						if ok, _ := env.Event.CheckSignature(); !ok {
							InfoLogger.Printf("{%s} bad signature on %s\n", r.URL, env.Event.ID)
							continue
						}
					}

					// dispatch this to the internal .events channel of the subscription
					subscription.dispatchEvent(&env.Event)
				}
			case *EOSEEnvelope:
				if subscription, ok := r.Subscriptions.Load(subIdToSerial(string(*env))); ok {
					subscription.dispatchEose()
				}
			case *ClosedEnvelope:
				if subscription, ok := r.Subscriptions.Load(subIdToSerial(env.SubscriptionID)); ok {
					subscription.handleClosed(env.Reason)
				}
			case *CountEnvelope:
				if subscription, ok := r.Subscriptions.Load(subIdToSerial(env.SubscriptionID)); ok && env.Count != nil && subscription.countResult != nil {
					subscription.countResult <- *env.Count
				}
			case *OKEnvelope:
				if okCallback, exist := r.okCallbacks.Load(env.EventID); exist {
					okCallback(env.OK, env.Reason)
				} else {
					InfoLogger.Printf("{%s} got an unexpected OK message for event %s", r.URL, env.EventID)
				}
			}
		}
	}()

	return nil
}

// Write queues a message to be sent to the relay.
func (r *Relay) Write(msg []byte) <-chan error {
	ch := make(chan error)
	select {
	case r.writeQueue <- writeRequest{msg: msg, answer: ch}:
	case <-r.connectionContext.Done():
		go func() { ch <- fmt.Errorf("connection closed") }()
	}
	return ch
}

// Publish sends an "EVENT" command to the relay r as in NIP-01 and waits for an OK response.
func (r *Relay) Publish(ctx context.Context, event Event) error {
	return r.publish(ctx, event.ID, &EventEnvelope{Event: event})
}

// Auth sends an "AUTH" command client->relay as in NIP-42 and waits for an OK response.
func (r *Relay) Auth(ctx context.Context, sign func(event *Event) error) error {
	authEvent := Event{
		CreatedAt: Now(),
		Kind:      KindClientAuthentication,
		Tags: Tags{
			Tag{"relay", r.URL},
			Tag{"challenge", r.challenge},
		},
		Content: "",
	}
	if err := sign(&authEvent); err != nil {
		return fmt.Errorf("error signing auth event: %w", err)
	}

	return r.publish(ctx, authEvent.ID, &AuthEnvelope{Event: authEvent})
}

// publish can be used both for EVENT and for AUTH
func (r *Relay) publish(ctx context.Context, id string, env Envelope) error {
	var err error
	var cancel context.CancelFunc

	if _, ok := ctx.Deadline(); !ok {
		// if no timeout is set, force it to 7 seconds
		ctx, cancel = context.WithTimeoutCause(ctx, 7*time.Second, fmt.Errorf("given up waiting for an OK"))
		defer cancel()
	} else {
		// otherwise make the context cancellable so we can stop everything upon receiving an "OK"
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}

	// listen for an OK callback
	gotOk := false
	r.okCallbacks.Store(id, func(ok bool, reason string) {
		gotOk = true
		if !ok {
			err = fmt.Errorf("msg: %s", reason)
		}
		cancel()
	})
	defer r.okCallbacks.Delete(id)

	// publish event
	envb, _ := env.MarshalJSON()
	if err := <-r.Write(envb); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			// this will be called when we get an OK or when the context has been canceled
			if gotOk {
				return err
			}
			return ctx.Err()
		case <-r.connectionContext.Done():
			// this is caused when we lose connectivity
			return err
		}
	}
}

// Subscribe sends a "REQ" command to the relay r as in NIP-01.
// Events are returned through the channel sub.Events.
// The subscription is closed when context ctx is cancelled ("CLOSE" in NIP-01).
//
// Remember to cancel subscriptions, either by calling `.Unsub()` on them or ensuring their `context.Context` will be canceled at some point.
// Failure to do that will result in a huge number of halted goroutines being created.
func (r *Relay) Subscribe(ctx context.Context, filters Filters, opts ...SubscriptionOption) (*Subscription, error) {
	sub := r.PrepareSubscription(ctx, filters, opts...)

	if r.Connection == nil {
		return nil, fmt.Errorf("not connected to %s", r.URL)
	}

	if err := sub.Fire(); err != nil {
		return nil, fmt.Errorf("couldn't subscribe to %v at %s: %w", filters, r.URL, err)
	}

	return sub, nil
}

// PrepareSubscription creates a subscription, but doesn't fire it.
//
// Remember to cancel subscriptions, either by calling `.Unsub()` on them or ensuring their `context.Context` will be canceled at some point.
// Failure to do that will result in a huge number of halted goroutines being created.
func (r *Relay) PrepareSubscription(ctx context.Context, filters Filters, opts ...SubscriptionOption) *Subscription {
	current := subscriptionIDCounter.Add(1)
	ctx, cancel := context.WithCancel(ctx)

	sub := &Subscription{
		Relay:             r,
		Context:           ctx,
		cancel:            cancel,
		counter:           current,
		Events:            make(chan *Event),
		EndOfStoredEvents: make(chan struct{}, 1),
		ClosedReason:      make(chan string, 1),
		Filters:           filters,
		match:             filters.Match,
	}

	label := ""
	for _, opt := range opts {
		switch o := opt.(type) {
		case WithLabel:
			label = string(o)
		}
	}

	// subscription id calculation
	buf := subIdPool.Get().([]byte)[:0]
	buf = strconv.AppendInt(buf, sub.counter, 10)
	buf = append(buf, ':')
	buf = append(buf, label...)
	defer subIdPool.Put(buf)
	sub.id = string(buf)

	// we track subscriptions only by their counter, no need for the full id
	r.Subscriptions.Store(int64(sub.counter), sub)

	// start handling events, eose, unsub etc:
	go sub.start()

	return sub
}

func (r *Relay) QueryEvents(ctx context.Context, filter Filter) (chan *Event, error) {
	sub, err := r.Subscribe(ctx, Filters{filter})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-sub.ClosedReason:
			case <-sub.EndOfStoredEvents:
			case <-ctx.Done():
			case <-r.Context().Done():
			}
			sub.Unsub()
			return
		}
	}()

	return sub.Events, nil
}

func (r *Relay) QuerySync(ctx context.Context, filter Filter) ([]*Event, error) {
	if _, ok := ctx.Deadline(); !ok {
		// if no timeout is set, force it to 7 seconds
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 7*time.Second)
		defer cancel()
	}

	events := make([]*Event, 0, max(filter.Limit, 250))
	ch, err := r.QueryEvents(ctx, filter)
	if err != nil {
		return nil, err
	}

	for evt := range ch {
		events = append(events, evt)
	}

	return events, nil
}

func (r *Relay) Count(ctx context.Context, filters Filters, opts ...SubscriptionOption) (int64, error) {
	sub := r.PrepareSubscription(ctx, filters, opts...)
	sub.countResult = make(chan int64)

	if err := sub.Fire(); err != nil {
		return 0, err
	}

	defer sub.Unsub()

	if _, ok := ctx.Deadline(); !ok {
		// if no timeout is set, force it to 7 seconds
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 7*time.Second)
		defer cancel()
	}

	for {
		select {
		case count := <-sub.countResult:
			return count, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
}

func (r *Relay) Close() error {
	r.closeMutex.Lock()
	defer r.closeMutex.Unlock()

	if r.connectionContextCancel == nil {
		return fmt.Errorf("relay already closed")
	}
	r.connectionContextCancel()
	r.connectionContextCancel = nil

	if r.Connection == nil {
		return fmt.Errorf("relay not connected")
	}

	err := r.Connection.Close()
	if err != nil {
		return err
	}

	return nil
}

var subIdPool = sync.Pool{
	New: func() any { return make([]byte, 0, 15) },
}
