// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cism_test

import (
	"errors"
	"github.com/sebuckler/cism"
	"testing"
)

func TestMachine_Start(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should err when states missing":          shouldErrStartStatesMissing,
		"should err when state not defined":       shouldErrStartStateNotDefined,
		"should err when machine stopped":         shouldErrStartMachineStopped,
		"should err when machine already started": shouldErrStartMachineStarted,
		"should be nil when successful":           shouldSucceedStart,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func TestMachine_Send(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should err when machine not started":                shouldErrSendMachineNotStarted,
		"should err when machine stopped":                    shouldErrSendMachineStopped,
		"should err when no transition for event":            shouldErrSendNoTran,
		"should be nil when transition successful":           shouldSucceedSend,
		"should stop machine when transition final":          shouldStopMachineFinalTran,
		"should handle failed transition when guard fails":   shouldHandleTranGuardFail,
		"should handle success transition when guard passes": shouldHandleTranGuardPass,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func TestMachine_Reset(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should err when machine not stopped": shouldErrResetMachineNotStopped,
		"should be nil when reset successful": shouldSucceedReset,
		"should start after reset":            shouldResetThenStart,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func TestMachine_Current(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should be correct state after transition": shouldSucceedCurrent,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func TestMachine_History(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should be empty when no transitions occurred": shouldBeEmptyHistoryNoTran,
		"should have correct transition history":       shouldSucceedHistory,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func shouldErrStartStatesMissing(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{}
	var machineErr *cism.ErrMissingStates

	if err := machine.Start(state); err == nil || err.Error() == "" || !errors.As(err, &machineErr) {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldErrStartStateNotDefined(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{cism.State(2): nil}}
	var machineErr *cism.ErrStateNotDefined

	if err := machine.Start(state); err == nil || err.Error() == "" || !errors.As(err, &machineErr) {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldErrStartMachineStopped(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: nil}}
	machine.Stop()
	var machineErr *cism.ErrMachineStopped

	if err := machine.Start(state); err == nil || err.Error() == "" || !errors.As(err, &machineErr) {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldErrStartMachineStarted(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: nil}}
	startErr := machine.Start(state)
	var machineErr *cism.ErrMachineStarted

	if err := machine.Start(state); err == nil || err.Error() == "" || !errors.As(err, &machineErr) || startErr != nil {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldSucceedStart(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: nil}}

	if err := machine.Start(state); err != nil {
		t.Fail()
		t.Logf("%s: errored", name)
	}
}

func shouldErrSendMachineNotStarted(t *testing.T, name string) {
	event := cism.Event(1)
	machine := &cism.Machine{}
	var errMachine *cism.ErrMachineNotStarted

	if err := machine.Send(event); err == nil || err.Error() == "" || !errors.As(err, &errMachine) {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldErrSendMachineStopped(t *testing.T, name string) {
	event := cism.Event(1)
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: nil}}}
	startErr := machine.Start(state)
	machine.Stop()
	var errMachine *cism.ErrMachineStopped

	if err := machine.Send(event); err == nil || !errors.As(err, &errMachine) || startErr != nil {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldErrSendNoTran(t *testing.T, name string) {
	event := cism.Event(1)
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: nil}}}
	startErr := machine.Start(state)
	var errMachine *cism.ErrMissingTransition

	if err := machine.Send(event); err == nil || err.Error() == "" || !errors.As(err, &errMachine) || startErr != nil {
		t.Fail()
		t.Logf("%s: did not error correctly", name)
	}
}

func shouldSucceedSend(t *testing.T, name string) {
	event := cism.Event(1)
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{}}}}
	startErr := machine.Start(state)

	if err := machine.Send(event); err != nil || startErr != nil {
		t.Fail()
		t.Logf("%s: errored", name)
	}
}

func shouldStopMachineFinalTran(t *testing.T, name string) {
	event := cism.Event(1)
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{IsFinal: true}}}}
	startErr := machine.Start(state)
	sendErr := machine.Send(event)

	if err := machine.Send(event); err == nil || sendErr != nil || startErr != nil {
		t.Fail()
		t.Logf("%s: machine not stopped", name)
	}
}

func shouldHandleTranGuardFail(t *testing.T, name string) {
	event := cism.Event(1)
	state := cism.State(1)
	state2 := cism.State(2)
	handled := false
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{
		Guard: func(s cism.State, e cism.Event) bool {
			return false
		},
		OnFail: func(s cism.State, e cism.Event) {
			handled = true
		},
		To: state2,
	}}}}
	startErr := machine.Start(state)

	if err := machine.Send(event); err != nil || startErr != nil || !handled || machine.Current() != state {
		t.Fail()
		t.Logf("%s: failed transition not handled", name)
	}
}

func shouldHandleTranGuardPass(t *testing.T, name string) {
	event := cism.Event(1)
	state := cism.State(1)
	state2 := cism.State(2)
	handled := false
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{
		Guard: func(s cism.State, e cism.Event) bool {
			return true
		},
		OnSuccess: func(s cism.State, e cism.Event) {
			handled = true
		},
		To: state2,
	}}}}
	startErr := machine.Start(state)

	if err := machine.Send(event); err != nil || startErr != nil || !handled || machine.Current() != state2 {
		t.Fail()
		t.Logf("%s: successful transition not handled", name)
	}
}

func shouldErrResetMachineNotStopped(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: nil}}
	startErr := machine.Start(state)
	var machineErr *cism.ErrMachineNotStopped

	if err := machine.Reset(); err == nil || err.Error() == "" || !errors.As(err, &machineErr) || startErr != nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldSucceedReset(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: nil}}
	startErr := machine.Start(state)
	machine.Stop()

	if err := machine.Reset(); err != nil || startErr != nil {
		t.Fail()
		t.Logf("%s: errored", name)
	}
}

func shouldResetThenStart(t *testing.T, name string) {
	state := cism.State(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: nil}}
	startErr := machine.Start(state)
	machine.Stop()
	resetErr := machine.Reset()

	if err := machine.Start(state); err != nil || startErr != nil || resetErr != nil {
		t.Fail()
		t.Logf("%s: errored", name)
	}
}

func shouldSucceedCurrent(t *testing.T, name string) {
	state := cism.State(1)
	state2 := cism.State(2)
	event := cism.Event(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{To: state2}}}}
	startErr := machine.Start(state)
	sendErr := machine.Send(event)

	if machine.Current() != state2 || startErr != nil || sendErr != nil {
		t.Fail()
		t.Logf("%s: did not match", name)
	}
}

func shouldBeEmptyHistoryNoTran(t *testing.T, name string) {
	state := cism.State(1)
	state2 := cism.State(2)
	event := cism.Event(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{To: state2}}}}
	startErr := machine.Start(state)

	if len(machine.History()) > 0 || startErr != nil {
		t.Fail()
		t.Logf("%s: not empty", name)
	}
}

func shouldSucceedHistory(t *testing.T, name string) {
	state := cism.State(1)
	state2 := cism.State(2)
	event := cism.Event(1)
	machine := &cism.Machine{States: cism.StateTransitionTable{state: {event: &cism.Transition{To: state2}}}}
	startErr := machine.Start(state)
	sendErr := machine.Send(event)

	if len(machine.History()) != 1 || startErr != nil || sendErr != nil {
		t.Fail()
		t.Logf("%s: no history", name)
	}
}
