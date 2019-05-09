package event

import (
	"fmt"
	"sync"

	"github.com/neovim/go-client/nvim"
)

type (
	Type         int
	CallbackFunc func() error

	handler struct {
		vim          *nvim.Nvim
		mux          sync.RWMutex
		callbacksMap map[Type][]CallbackFunc
	}
)

const (
	Group = "trans-event"

	TypeMoveEvent Type = iota
)

var (
	defaultHandler *handler
)

func RegisterHandler(vim *nvim.Nvim) {
	defaultHandler = &handler{
		vim:          vim,
		callbacksMap: map[Type][]CallbackFunc{},
	}
}

func Callback(t Type) CallbackFunc {
	return defaultHandler.Callback(t)
}

func On(t Type, cb CallbackFunc) {
	defaultHandler.On(t, cb)
}

func (h *handler) Callback(t Type) CallbackFunc {
	return func() error {
		h.mux.Lock()
		defer h.mux.Unlock()

		cbs := h.callbacksMap[t]
		if len(cbs) == 0 {
			return nil
		}

		cb := cbs[0]
		h.callbacksMap[t] = cbs[1:]

		go func() {
			if err := cb(); err != nil {
				h.vim.WriteErr(fmt.Sprintf("failed to callback function on %v: %v\n", t, err))
			}
		}()
		return nil
	}
}

func (h *handler) On(t Type, cb CallbackFunc) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.callbacksMap[t] = append(h.callbacksMap[t], cb)
}
