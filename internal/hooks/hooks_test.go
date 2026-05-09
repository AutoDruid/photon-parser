package hooks_test

import (
	"testing"

	"AutoDruid/photon-parser/internal/hooks"
	v16 "AutoDruid/photon-parser/internal/parameters/v16"
	v18 "AutoDruid/photon-parser/internal/parameters/v18"
	"AutoDruid/photon-parser/internal/types"
)

func TestNewHooks(t *testing.T) {

	h := hooks.NewHooks[v18.Parameter]()
	if h == nil {
		t.Fatalf("NewHooks: expected non-nil hooks")
	}

	v16h := hooks.NewHooks[v16.Parameter]()
	if v16h == nil {
		t.Fatalf("NewHooks: expected non-nil hooks")
	}
}

func TestAsyncHooks_ReusesSameChannel(t *testing.T) {
	t.Parallel()

	opts := types.HookOptions{Size: 4}
	h := &hooks.Hooks[v18.Parameter]{}

	s1 := h.OnSessionAsync(opts)
	s2 := h.OnSessionAsync(opts)
	if s1 != s2 {
		t.Fatalf("OnSessionAsync: expected same channel instance on second call")
	}
	if cap(s1) != int(opts.Size) {
		t.Fatalf("OnSessionAsync: got cap %d, want %d", cap(s1), opts.Size)
	}

	c1 := h.OnCommandAsync(opts)
	c2 := h.OnCommandAsync(opts)
	if c1 != c2 {
		t.Fatalf("OnCommandAsync: expected same channel instance on second call")
	}
	if cap(c1) != int(opts.Size) {
		t.Fatalf("OnCommandAsync: got cap %d, want %d", cap(c1), opts.Size)
	}

	p1 := h.OnParameterAsync(opts)
	p2 := h.OnParameterAsync(opts)
	if p1 != p2 {
		t.Fatalf("OnParameterAsync: expected same channel instance on second call")
	}
	if cap(p1) != int(opts.Size) {
		t.Fatalf("OnParameterAsync: got cap %d, want %d", cap(p1), opts.Size)
	}
}

func TestAsyncHooks_BufferSizeFromFirstCall(t *testing.T) {
	t.Parallel()

	h := &hooks.Hooks[v18.Parameter]{}
	first := types.HookOptions{Size: 2}
	second := types.HookOptions{Size: 99}

	ch := h.OnSessionAsync(first)
	_ = h.OnSessionAsync(second)
	if cap(ch) != 2 {
		t.Fatalf("expected buffer from first OnSessionAsync (2), got %d", cap(ch))
	}
}

func TestCloseAsyncHooks_ReceiveClosedZeroValue(t *testing.T) {
	t.Parallel()

	h := &hooks.Hooks[v18.Parameter]{}
	opts := types.HookOptions{Size: 1}

	sessCh := h.OnSessionAsync(opts)
	cmdCh := h.OnCommandAsync(opts)
	paramCh := h.OnParameterAsync(opts)

	h.CloseAsyncHooks()

	if _, ok := <-sessCh; ok {
		t.Fatalf("session channel: expected receive on closed channel with ok=false")
	}
	if _, ok := <-cmdCh; ok {
		t.Fatalf("command channel: expected receive on closed channel with ok=false")
	}
	if _, ok := <-paramCh; ok {
		t.Fatalf("parameter channel: expected receive on closed channel with ok=false")
	}
}

func TestCloseAsyncHooks_AllowsNewChannels(t *testing.T) {
	t.Parallel()

	opts := types.HookOptions{Size: 1}
	h := &hooks.Hooks[v18.Parameter]{}

	s1 := h.OnSessionAsync(opts)
	h.CloseAsyncHooks()
	s2 := h.OnSessionAsync(opts)
	if s1 == s2 {
		t.Fatalf("expected new session channel after CloseAsyncHooks")
	}
	if _, ok := <-s1; ok {
		t.Fatalf("old session channel should be closed")
	}

	c1 := h.OnCommandAsync(opts)
	h.CloseAsyncHooks()
	c2 := h.OnCommandAsync(opts)
	if c1 == c2 {
		t.Fatalf("expected new command channel after CloseAsyncHooks")
	}

	p1 := h.OnParameterAsync(opts)
	h.CloseAsyncHooks()
	p2 := h.OnParameterAsync(opts)
	if p1 == p2 {
		t.Fatalf("expected new parameter channel after CloseAsyncHooks")
	}
}
