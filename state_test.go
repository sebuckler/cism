// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cism_test

import (
	"github.com/sebuckler/cism"
	"testing"
)

func TestStateTransitionTable_GetTransition(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should be nil when table is empty":                shouldBeNilTranEmptyTable,
		"should be nil when state is missing":              shouldBeNilTranMissingState,
		"should exist when state and event match in table": shouldExistTranOnMatch,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func TestStateTransitionTable_GetEventsForState(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should be empty when state is missing":    shouldBeEmptyEventsMissingState,
		"should be empty when state has no events": shouldBeEmptyEventsNoEvents,
		"should exist when state has events":       shouldExistEventsOnMatch,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func TestStateTransitionTable_GetStatesForEvent(t *testing.T) {
	testCases := map[string]func(t *testing.T, name string){
		"should be empty when table is empty":      shouldBeEmptyStatesEmptyTable,
		"should be empty when event has no states": shouldBeEmptyStatesNoStates,
		"should exist when event has states":       shouldExistStatesOnMatch,
	}

	for name, test := range testCases {
		test(t, name)
	}
}

func shouldBeNilTranEmptyTable(t *testing.T, name string) {
	stt := cism.StateTransitionTable{}

	if stt.GetTransition(cism.State(1), cism.Event(1)) != nil {
		t.Fail()
		t.Logf("%s: result not nil", name)
	}
}

func shouldBeNilTranMissingState(t *testing.T, name string) {
	state := cism.State(1)
	event := cism.Event(1)
	stt := cism.StateTransitionTable{cism.State(2): {event: &cism.Transition{}}}

	if stt.GetTransition(state, event) != nil {
		t.Fail()
		t.Logf("%s: result not nil", name)
	}
}

func shouldExistTranOnMatch(t *testing.T, name string) {
	state := cism.State(1)
	event := cism.Event(1)
	tran := &cism.Transition{}
	stt := cism.StateTransitionTable{state: {event: tran}}

	if stt.GetTransition(state, event) != tran {
		t.Fail()
		t.Logf("%s: incorrect transaction returned", name)
	}
}

func shouldBeEmptyEventsMissingState(t *testing.T, name string) {
	state := cism.State(1)
	stt := cism.StateTransitionTable{cism.State(2): {cism.Event(1): nil}}

	if len(stt.GetEventsForState(state)) > 0 {
		t.Fail()
		t.Logf("%s: result not empty", name)
	}
}

func shouldBeEmptyEventsNoEvents(t *testing.T, name string) {
	state := cism.State(1)
	stt := cism.StateTransitionTable{state: {}}

	if len(stt.GetEventsForState(state)) > 0 {
		t.Fail()
		t.Logf("%s: result not empty", name)
	}
}

func shouldExistEventsOnMatch(t *testing.T, name string) {
	state := cism.State(1)
	event := cism.Event(1)
	stt := cism.StateTransitionTable{state: {event: nil}}

	if len(stt.GetEventsForState(state)) != 1 {
		t.Fail()
		t.Logf("%s: result empty", name)
	}
}

func shouldBeEmptyStatesEmptyTable(t *testing.T, name string) {
	event := cism.Event(1)
	stt := cism.StateTransitionTable{}

	if len(stt.GetStatesForEvent(event)) > 0 {
		t.Fail()
		t.Logf("%s: result not empty", name)
	}
}

func shouldBeEmptyStatesNoStates(t *testing.T, name string) {
	state := cism.State(1)
	event := cism.Event(1)
	stt := cism.StateTransitionTable{state: {cism.Event(2): nil}}

	if len(stt.GetStatesForEvent(event)) > 0 {
		t.Fail()
		t.Logf("%s: result not empty", name)
	}
}

func shouldExistStatesOnMatch(t *testing.T, name string) {
	state := cism.State(1)
	event := cism.Event(1)
	stt := cism.StateTransitionTable{state: {event: nil}}

	if len(stt.GetStatesForEvent(event)) != 1 {
		t.Fail()
		t.Logf("%s: result empty", name)
	}
}
