package streamer

import (
	"context"
)

type (
	// A StreamManager accepts one or more Connections to be managed
	//
	// When registered Connections are kept open by periodically inspecting them and reconnecting them if needed.
	StreamManager struct {
		application   context.Context
		registrations map[string]*Reader
	}
)

// NewStreamManager returns a new instance of StreamManager reference.
//
// Note that the provided context should be used to terminate all activity in the registered Connectors
// ideally this is the context with spans the application lifecycle such as the one provided by Fx
func NewStreamManager(application context.Context) *StreamManager {
	m := &StreamManager{
		application:   application,
		registrations: make(map[string]*Reader),
	}
	return m
}

// EnsureConnected scans all streamers and reconnects them if needed.
// This method should be called periodically e.g. by pairing it with the periodic package
//
//     fx.Invoke(
//			func(l fx.Lifecycle, scheduler *periodic.Service, sm *client.StreamManager) {
//				l.Append(fx.Hook{
//					OnStart: func(ctx context.Context) error {
//						scheduler.StartPeriodicTask(func() {
//							sm.EnsureConnected(func(_ error){})
//						}, "connect-streams", 5*time.Second, 0)
//						return nil
//					},
//					OnStop: func(_ context.Context) error {
//						return nil
//					},
//				})
//			},
//		),
//
func (sm *StreamManager) EnsureConnected(onError func(error)) {
	for _, r := range sm.registrations {
		err := sm.ensureStart(r)
		if err != nil {
			onError(err)
		}
	}
}

func (sm *StreamManager) ensureStart(r *Reader) error {
	if !r.Active() {
		return r.Start()
	}
	return nil
}

// Register registers a Connector and StreamHandler with the manager
func (sm *StreamManager) Register(c Connector, h StreamHandler) {
	sm.RegisterFunc(c, h.Handle)
}

// RegisterFunc registers a Connector and StreamHandlerFunc with the manager
func (sm *StreamManager) RegisterFunc(c Connector, h StreamHandlerFunc) {
	connectorID := c.ID()
	_, exists := sm.registrations[connectorID]
	if !exists {
		sm.registrations[connectorID] = NewReader(sm.application, c, h)
	}
}
