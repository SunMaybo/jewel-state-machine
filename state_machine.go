package jewel_state_machine

import (
	"sync"
	"errors"
)

type CallBackType int

const (
	CallBackBeforeEvent      CallBackType = iota + 1
	CallBackEnterEvent
	CallBackAfterEvent
	CallBackThrowEvent
	CallBackAfterReturnEvent
)

func (ct CallBackType) String() string {
	switch ct {
	case CallBackAfterEvent:
		return "CALLBACK_AFTER_EVENT"
	case CallBackEnterEvent:
		return "CALLBACK_ENTER_EVENT"
	case CallBackBeforeEvent:
		return "CALLBACK_BEFORE_EVENT"
	case CallBackThrowEvent:
		return "CALLBACK_THROW_EVENT"
	}
	return "NO_EVENT"
}

type StateMachine struct {
	current   string
	events    []EventDesc
	stateMu   sync.RWMutex
	callbacks map[CallBackType]Callback
}

type EventDesc struct {
	Name string
	Src  []string
}

type Callback func(event *Event) error

func (cb Callback) Do(event *Event) error {
	return cb(event)
}

func New(current string, events []EventDesc, callbacks map[CallBackType]Callback) *StateMachine {
	sm := &StateMachine{
		callbacks: callbacks,
		current:   current,
		events:    events,
	}
	return sm
}

func (sm *StateMachine) Transaction(target string, args ...interface{}) error {
	//前置事件处理
	if v, ok := sm.callbacks[CallBackBeforeEvent]; ok {
		v(&Event{
			Src:   sm.current,
			Event: target,
			Args:  args,
		})
	}
	//
	flag := false
	var err error
	sm.stateMu.Lock()
loop:
	for _, event := range sm.events {
		if event.Name == target {
			for _, src := range event.Src {
				if src == sm.current {

					flag = true
					if v, ok := sm.callbacks[CallBackEnterEvent]; ok {
						err = v(&Event{
							Src:   sm.current,
							Event: target,
							Args:  args,
						})
					}
					if err == nil {
						sm.current = target
					}
					break loop

				}
			}
		}
	}

	sm.stateMu.Unlock()

	if !flag {
		err = errors.New("current state transaction to target state failed ")
	}

	if err != nil {
		if v, ok := sm.callbacks[CallBackThrowEvent]; ok {
			v(&Event{
				Src:   sm.current,
				Event: target,
				Args:  args,
				Err:   err,
			})
		}
	} else {
		if v, ok := sm.callbacks[CallBackAfterEvent]; ok {
			defer v(&Event{
				Src:   sm.current,
				Event: target,
				Args:  args,
				Err:   err,
			})
		}
	}

	if v, ok := sm.callbacks[CallBackAfterReturnEvent]; ok {
		defer v(&Event{
			Src:   sm.current,
			Event: target,
			Args:  args,
			Err:   err,
		})
	}
	return err

}

func (sm *StateMachine) Current() string {
	sm.stateMu.RLock()
	defer sm.stateMu.RUnlock()
	return sm.current
}

func (sm *StateMachine) Is(state string) bool {
	sm.stateMu.RLock()
	defer sm.stateMu.RUnlock()
	return state == sm.current
}
func (sm *StateMachine) SetState(state string) {
	sm.stateMu.Lock()
	defer sm.stateMu.Unlock()
	sm.current = state
	return
}
