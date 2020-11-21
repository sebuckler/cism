// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cism

/*
ErrMissingStates represents an error when a machine is attempting to be started
without a state transition table defined. It satisfies the Error interface.
*/
type ErrMissingStates struct {
	msg string
}

/*
Error returns the error message assigned at struct creation.
*/
func (e *ErrMissingStates) Error() string {
	return e.msg
}

/*
ErrStateNotDefined represents an error when a machine is attempting to be
started with a state that does not exist in the state transition table. It
satisfies the Error interface.
*/
type ErrStateNotDefined struct {
	State State
	msg   string
}

/*
Error returns the error message assigned at struct creation.
*/
func (e *ErrStateNotDefined) Error() string {
	return e.msg
}

/*
ErrMachineStopped represents an error when a machine is in the done state and
an attempt is made to start it or send an event to it. It satisfies the Error
interface.
*/
type ErrMachineStopped struct {
	FinalEvent *Event
	msg        string
}

/*
Error returns the error message assigned at struct creation.
*/
func (e *ErrMachineStopped) Error() string {
	return e.msg
}

/*
ErrMachineStarted represents an error when a machine is attempting to be
started when it has already been started. It satisfies the Error interface.
*/
type ErrMachineStarted struct {
	msg string
}

/*
Error returns the error message assigned at struct creation.
*/
func (e *ErrMachineStarted) Error() string {
	return e.msg
}

/*
ErrMissingTransition represents an error when a machine is attempting to have
an event sent to it and the event has no transition for the current state. It
satisfies the Error interface.
*/
type ErrMissingTransition struct {
	State State
	Event Event
	msg   string
}

func (e *ErrMissingTransition) Error() string {
	return e.msg
}

/*
ErrMachineNotStarted represents an error when a machine is attempting to have
an event sent to it and the machine has not been started. It satisfies the
Error interface.
*/
type ErrMachineNotStarted struct {
	msg string
}

/*
Error returns the error message assigned at struct creation.
*/
func (e *ErrMachineNotStarted) Error() string {
	return e.msg
}

/*
ErrMachineNotStopped represents an error when a machine is in the done state
and an attempt is made to reset it. It satisfies the Error interface.
*/
type ErrMachineNotStopped struct {
	msg string
}

/*
Error returns the error message assigned at struct creation.
*/
func (e *ErrMachineNotStopped) Error() string {
	return e.msg
}
