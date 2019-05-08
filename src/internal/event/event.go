package event

import "sync"

type (
	Type    int
	eventFn func()
)

const (
	Group = "trans-event"

	TypeMoveEvent Type = iota
)

var (
	mux         sync.RWMutex
	eventFnsMap = map[Type][]eventFn{}
)

func init() {
	eventFnsMap[TypeMoveEvent] = nil
}

func Callback(t Type) func() error {
	return func() error {
		mux.Lock()
		defer mux.Unlock()

		fns := eventFnsMap[t]
		if len(fns) == 0 {
			return nil
		}

		go fns[0]()

		eventFnsMap[t] = fns[1:]

		return nil
	}
}

func On(t Type, fn eventFn) {
	mux.Lock()
	defer mux.Unlock()

	eventFnsMap[t] = append(eventFnsMap[t], fn)
}

func (t Type) String() string {
	return string(t)
}
