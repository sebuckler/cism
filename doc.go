// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

/*
Package cism allows for using state machines defined with a state transition
table.

Basic Operations

Define states and events to be used as the keys of a state transition table.
States represent the current status of a status or program. Events are actions
or inputs that attempt to trigger state changes with a state machine. Create
transitions for each event for each state.

Transitions represent the lifecycle of a state change within a state machine.
Lifecycle hooks are exposed on transition definitions to allow for inverting
control of state-based logic to be called from the state machine after state
change attempts occur.

After defining the needed states, events, and transitions for setting up the
state transition table, create a state machine. The state machine takes the
state transition table as its instructions for how to handle incoming events
for the current state.

A state machine must be explicitly started to begin accepting events. Then, the
state machine can be stopped, reset, and started again, if needed. The current
state can be viewed for the purposes of logging or other introspections. A
history log is also kept that maintains past states and the events that
triggered the transitions.

Example code:

	package main

	import (
		"fmt"
		"github.com/sebuckler/cism"
	)

	const (
		Begin cism.State = iota
		Middle
		End
	)

	const (
		SetupDone cism.Event = iota
		WorkComplete
	)

	func main() {
		workReallyComplete := false
		stt := cism.StateTransitionTable{
			Begin: {
				SetupDone: &cism.Transition{
					OnSuccess: func(s cism.State, e cism.Event) {
						fmt.Println("left 'Begin' state and entered 'Middle' state")
					},
					To: Middle,
				},
			},
			Middle: {
				WorkComplete: &cism.Transition{
					Guard: func(s cism.State, e cism.Event) bool {
						return workReallyComplete
					},
					IsFinal: true,
					OnFail: func(s cism.State, e cism.Event) {
						workReallyComplete = true
					},
					OnSuccess: func(s cism.State, e cism.Event) {
						fmt.Println("left 'Middle' state and entered 'End' state")
					},
					To: End,
				},
			},
			End: {}, // done
		}
		machine := &cism.Machine{
			States: stt,
		}

		machine.Start(Begin)
		machine.Send(SetupDone)
		machine.Send(WorkComplete) // fails transition
		machine.Send(WorkComplete) // succeeds this time
		machine.Stop()
		machine.Reset()
	}
*/
package cism
