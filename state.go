// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cism

/*
State represents values to be used as state keys in a state transition table.
*/
type State int

/*
Event represents values to be used as event keys for states in a state
transition table.
*/
type Event int

/*
Transition is the context and lifecycle of a state change for an event in the
current state.
*/
type Transition struct {
	Guard     func(s State, e Event) bool // Lifecycle hook for allowing or blocking state change
	IsFinal   bool                        // Triggers machine done state if true
	OnFail    func(s State, e Event)      // Lifecycle hook for when Guard blocks state change
	OnSuccess func(s State, e Event)      // Lifecycle hook for when Guard allows state change
	To        State                       // State to transition to if Guard allows state change
}

/*
StateTransitionTable is a table mapping transitions to events for each state.
*/
type StateTransitionTable map[State]map[Event]*Transition

/*
GetTransition attempts to return a transition for a given state and event.
If the state does not exist in the table, it will return nil.
*/
func (stt StateTransitionTable) GetTransition(s State, e Event) *Transition {
	if _, ok := stt[s]; !ok {
		return nil
	}

	return stt[s][e]
}

/*
GetEventsForState attempts to return a slice of events for a given state. If
the state does not exist in the table or no events exist for the given state,
it will return an empty slice.
*/
func (stt StateTransitionTable) GetEventsForState(s State) []Event {
	if _, ok := stt[s]; !ok {
		return nil
	}

	var events []Event

	for event, _ := range stt[s] {
		events = append(events, event)
	}

	return events
}

/*
GetStatesForEvent attempts to return a slice of states for a given event. If
the table is empty or no states have the given event, it will return an empty
slice.
*/
func (stt StateTransitionTable) GetStatesForEvent(e Event) []State {
	if len(stt) == 0 {
		return nil
	}

	var states []State

	for state, events := range stt {
		for event, _ := range events {
			if event == e {
				states = append(states, state)
			}
		}
	}

	return states
}
