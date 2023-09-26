package bus

import (
	"context"
	"sync"
	"testing"
)

type HandlerMock struct {
	I    int
	lock sync.Mutex
}

func (h *HandlerMock) Handle(ctx context.Context, hi any) {
	h.lock.Lock()
	h.I++
	h.lock.Unlock()
}

func TestBus_GetHandlers(t *testing.T) {
	b := NewBus[any]()
	h := &HandlerMock{}
	b.RegisterHandler("topic", h)
	if handlers, err := b.GetHandlers("topic"); err != nil || len(handlers) != 1 {
		t.Log("Expected to be registered only one handler")
		t.Fail()
	}
}

func TestBus_RegisterHandler(t *testing.T) {
	b := NewBus[any]()
	b.RegisterHandler("topic1", &HandlerMock{})
	b.RegisterHandler("topic1", &HandlerMock{})
	if handlers, err := b.GetHandlers("topic1"); err != nil || len(handlers) != 2 {
		t.Log("Expected to be registered two handlers")
		t.Fail()
	}
	handler := &HandlerMock{}
	b.RegisterHandler("topic2", handler)
	if handlers, err := b.GetHandlers("topic2"); err != nil || len(handlers) != 1 {
		t.Log("Expected to be registered only one handler")
		t.Fail()
	}
	if err := b.RegisterHandler("topic2", handler); err == nil {
		t.Log("Expected not to register duplicate handlers")
		t.Fail()
	}

}

func TestBus_UnregisterHandler(t *testing.T) {
	b := NewBus[any]()
	handler := &HandlerMock{}
	b.RegisterHandler("topic1", handler)
	if err := b.UnregisterHandler("topic1", handler); err != nil {
		t.Fail()
	}
	b.RegisterHandler("topic1", handler)
	if err := b.UnregisterHandler("topic1", &HandlerMock{}); err == nil {
		t.Log("Expected not to unregister the unregistered handler and returns an error")
		t.Fail()
	}
}

func TestBus_Broadcast(t *testing.T) {
	b := NewBus[any]()
	handler := &HandlerMock{I: 0, lock: sync.Mutex{}}
	b.RegisterHandler("topic", handler)
	b.Broadcast(context.Background(), "topic", struct{}{})
	if handler.I != 1 {
		t.Logf("Expected to call the handler exactly once but got %d", handler.I)
		t.Fail()
	}
}

func TestBus_BroadcastAsync(t *testing.T) {
	b := NewBus[any]()
	handler := &HandlerMock{I: 0, lock: sync.Mutex{}}
	b.RegisterHandler("topic", handler)
	for i := 0; i < 100; i++ {
		b.BroadcastAsync(context.Background(), "topic", struct{}{})
	}
	b.WaitAsync()
	if handler.I != 100 {
		t.Logf("Expected handler has been called 100 times but got %d", handler.I)
		t.Fail()
	}
}
