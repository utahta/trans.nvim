package event

import (
	"fmt"
	"sync"

	"github.com/neovim/go-client/nvim"
)

type (
	Type     int
	Callback func() error

	handler struct {
		vim       *nvim.Nvim
		mux       sync.RWMutex
		callbacks map[Type][]Callback
	}
)

const (
	Group = "trans-event"

	TypeCursorMoved Type = iota
)

var (
	defaultHandler     *handler
	defaultHandlerOnce sync.Once
)

func Init(vim *nvim.Nvim) {
	defaultHandlerOnce.Do(func() {
		defaultHandler = &handler{
			vim:       vim,
			callbacks: make(map[Type][]Callback),
		}
	})
}

func HandleFunc(t Type) func() error {
	return defaultHandler.HandleFunc(t)
}

func On(t Type, cb Callback) {
	defaultHandler.On(t, cb)
}

func (h *handler) HandleFunc(t Type) func() error {
	return func() error {
		h.mux.Lock()
		defer h.mux.Unlock()

		cbs := h.callbacks[t]
		if len(cbs) == 0 {
			return nil
		}

		cb := cbs[0]
		h.callbacks[t] = cbs[1:]

		go func() {
			if err := cb(); err != nil {
				h.vim.WriteErr(fmt.Sprintf("failed to callback function on %v: %v\n", t, err))
			}
		}()
		return nil
	}
}

func (h *handler) On(t Type, cb Callback) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.callbacks[t] = append(h.callbacks[t], cb)
}
