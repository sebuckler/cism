// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cism

/*
HistoryRecord represents a past state change and the event that caused it.
*/
type HistoryRecord struct {
	State State // state that was transitioned to
	Event Event // event that triggered the state change
}

/*
Machine is a state machine driven by a state transition table. All state
transitions are managed by the state machine.
*/
type Machine struct {
	States  StateTransitionTable // states and events the machine uses for transitions
	curr    State
	done    bool
	endevt  *Event
	hist    []HistoryRecord
	initial State
	started bool
}

/*
Start will start the machine at the given state. It will return an error if the
state transition table was not set on the machine. It will return an error if
the given state does not exist in the state transition table. It will return an
error if the machine has been stopped. It will return an error if the machine
has already been started.
*/
func (m *Machine) Start(s State) error {
	if len(m.States) == 0 {
		return &ErrMissingStates{"no states set"}
	}

	if _, ok := m.States[s]; !ok {
		return &ErrStateNotDefined{s, "start state not defined in states"}
	}

	if m.done {
		return &ErrMachineStopped{m.endevt, "machine is done and cannot be started"}
	}

	if m.started {
		return &ErrMachineStarted{"machine has already started"}
	}

	m.curr = s
	m.done = false
	m.initial = s
	m.started = true

	return nil
}

/*
Send will attempt to begin a state change based on the given event and the
current state. It will return an error if the machine has not been started. It
will return an error if the machine has been stopped. It will return an error
if no transition is defined for the given event and current state. The
transition lifecycle hooks will be invoked to determine if the machine can
complete the state change. If the transition guard fails, a failed state change
handler will be invoked. If the transition guard passes, a successful state
change handler will be invoked. If the transition is marked as final, the
machine will be stopped after the state change. If the state change succeeds,
the current history and triggering event will be pushed into a history log.
*/
func (m *Machine) Send(e Event) error {
	if !m.started {
		return &ErrMachineNotStarted{"machine has not started"}
	}

	if m.done {
		return &ErrMachineStopped{m.endevt, "machine is done and not accepting transitions"}
	}

	if tran := m.States.GetTransition(m.curr, e); tran != nil {
		m.transition(tran, e)
	} else {
		return &ErrMissingTransition{m.curr, e, "no transition found for event in current state"}
	}

	return nil
}

/*
Stop marks the machine as stopped and will accept no more state changes.
*/
func (m *Machine) Stop() {
	m.stop()
}

/*
Reset marks the machine as not stopped and not started. It will return an error
if the machine has not been stopped. The history log will be cleared on a
successful reset. The machine can be started again after it has been reset.
*/
func (m *Machine) Reset() error {
	if !m.done {
		return &ErrMachineNotStopped{"machine has not stopped"}
	}

	m.done = false
	m.endevt = nil
	m.hist = nil
	m.curr = m.initial
	m.started = false

	return nil
}

/*
Current returns the current state the machine is in.
*/
func (m *Machine) Current() State {
	return m.curr
}

/*
History returns a copy of the machine's state change history log.
*/
func (m *Machine) History() []HistoryRecord {
	cpyhist := make([]HistoryRecord, len(m.hist))

	copy(cpyhist, m.hist)

	return cpyhist
}

func (m *Machine) transition(tran *Transition, e Event) {
	currstate := m.curr

	if tran.Guard == nil || tran.Guard(currstate, e) {
		m.hist = append(m.hist, HistoryRecord{currstate, e})
		m.curr = tran.To

		if tran.OnSuccess != nil {
			tran.OnSuccess(currstate, e)
		}

		if tran.IsFinal {
			m.stop()
		}
	} else if tran.OnFail != nil {
		tran.OnFail(currstate, e)
	}
}

func (m *Machine) stop() {
	histlen := len(m.hist)

	if histlen > 0 {
		m.endevt = &(m.hist[histlen-1].Event)
	}

	m.done = true
}
