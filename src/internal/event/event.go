package event

import "sync"

type (
	Type          int
	eventCallback func() error
)

const (
	Group = "trans-event"

	TypeMoveEvent Type = iota
)

var (
	mux               sync.RWMutex
	eventCallbacksMap = map[Type][]eventCallback{}
)

func init() {
	eventCallbacksMap[TypeMoveEvent] = nil
}

func Callback(t Type) func() error {
	return func() error {
		mux.Lock()
		defer mux.Unlock()

		cbs := eventCallbacksMap[t]
		if len(cbs) == 0 {
			return nil
		}

		cb := cbs[0]
		eventCallbacksMap[t] = cbs[1:]

		go func() {
			if err := cb(); err != nil {
				panic(err)
			}
		}()
		return nil
	}
}

func On(t Type, cb eventCallback) {
	mux.Lock()
	defer mux.Unlock()

	eventCallbacksMap[t] = append(eventCallbacksMap[t], cb)
}

func (t Type) String() string {
	return string(t)
}
