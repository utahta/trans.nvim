package event

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

type (
	Callback func() error

	handler struct {
		vim       *nvim.Nvim
		mux       sync.RWMutex
		callbacks map[string]Callback

		group string
	}
)

var (
	defaultHandler     *handler
	defaultHandlerOnce sync.Once
)

const (
	handleFuncName = "TransEventHandle"
)

func Init(p *plugin.Plugin, group string) {
	defaultHandlerOnce.Do(func() {
		defaultHandler = &handler{
			vim:       p.Nvim,
			callbacks: make(map[string]Callback),

			group: group,
		}

		p.HandleFunction(&plugin.FunctionOptions{
			Name: handleFuncName,
		}, defaultHandler.handle)
	})
}

func Once(typ string, pattern string, cb Callback) {
	defaultHandler.once(typ, pattern, cb)
}

func (h *handler) handle(args []string) {
	if len(args) != 1 {
		h.vim.WritelnErr("event.handler.handle: invalid argument")
		return
	}
	id := args[0]

	h.mux.Lock()
	defer h.mux.Unlock()

	cb, ok := h.callbacks[id]
	if !ok {
		return
	}
	delete(h.callbacks, id)

	go func() {
		if err := cb(); err != nil {
			h.vim.WritelnErr(fmt.Sprintf("event.handler.handle: an error has occurred. err:%v", err))
		}
	}()
}

func (h *handler) buildHandleCommand(cb Callback) string {
	id := uuid.New().String()

	h.mux.Lock()
	h.callbacks[id] = cb
	h.mux.Unlock()

	return fmt.Sprintf(`call %s("%s")`, handleFuncName, id)
}

func (h *handler) once(typ string, pattern string, cb Callback) {
	h.on(typ, pattern, true, cb)
}

func (h *handler) on(typ string, pattern string, once bool, cb Callback) {
	var options []string
	if once {
		options = append(options, "++once")
	}
	cmd := fmt.Sprintf(`augroup %s | au %s %s %s %s | augroup END`,
		h.group,
		typ,
		pattern,
		strings.Join(options, " "),
		h.buildHandleCommand(cb),
	)
	if err := h.vim.Command(cmd); err != nil {
		h.vim.WritelnErr(fmt.Sprintf("event.handler.once: an error has occurred. err:%v", err))
	}
}
