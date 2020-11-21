# _CISM_

_CISM_ stands for _control-inverted state machine_.
This package is for adding the reliability of a state machine to an application.
States, events, and transitions are defined using a state transition table.
The inverted control nature of the state machine allows for the ability to hook into the transition lifecycle.
This means there is no need to ask the state machine what the current state is.
The state machine will call registered transition lifecycle functions when a transition is occurring.

## Installation

Use `go get` to install the latest version.

```
go get github.com/sebuckler/cism
```

Import `cism` like any other package.
```go
import "github.com/sebuckler/cism" 
```

## Usage

Define states, events, and transitions inside of a state transition table.
Then, create a machine using that table and send events to invoke state changes.

### State Transition Table

The state transition table holds all the available states a machine can transition to based on incoming events.
The states and events form a compound key with a transition that dictates the next state.
States don't have to have any events, in which case they are essentially a final state.
If an event is defined for a state, there _must_ be a transition defined, as well.

Create a state transition table struct.

```go
stt := cism.StateTransitionTable{}
```

#### State

States are representations of point-in-time instances of a system.
The simplest expression of a state is as an enumerated set.
States _must_ be unique.
If there are duplicate state values, the last one defined will be used in the state transition table.

Create states in the likeness of an `enum`.

```go
const (
	Begin cism.State = iota
	Middle
	End
)
```

Add the states to the state transition table.

```go
stt[Begin] = map[cism.Event]*cism.Transition{}
stt[Middle] = map[cism.Event]*cism.Transition{}
stt[End] = map[cism.Event]*cism.Transition{}
```

#### Event

Events are actions that are specific to a state and tell the machine to attempt a state change.
States can have multiple events defined, and events can be defined for multiple states.
If an event is defined on a state, it _must_ have a transition defined, as well.
The simplest expression of an event is as an enumerated set.
Events _must_ be unique.
If there are duplicate event values, the last one defined will be used in the state transition table.

Create events in the likeness of an `enum`.

```go
const (
	SetupDone cism.Event = iota
	WorkComplete
)
```

Add the events to the state transition table for the states where they can be triggered.

```go
stt[Begin][SetupDone] = &cism.Transition{}
stt[Middle][WorkComplete] = &cism.Transition{}
```

#### Transition

Transitions represent the attempt and lifecycle of a state change triggered by an event.
Each transition can define properties that may affect, or run as a result of, an attempted state change.

Create a transition struct from the `Begin` state to the `Middle` state, triggered by the `SetupDone` event.

```go
stt[Begin][SetupDone] = &cism.Transition{
    Guard: func(s cism.State, e cism.Event) bool {
        fmt.Printf("current state: %d\n", s)
        fmt.Printf("triggering event: %d\n", e)

        return true
    },
    IsFinal: false,
    OnFail: func(s cism.State, e cism.Event) {
        fmt.Println("failed to transition, no state change occurred")
    },
    OnSuccess: func(s cism.State, e cism.Event) {
        fmt.Println("transition succeeded, machine is in new state")
    },
    To: Middle,
}
```

 * `Guard` is a function that is used to determine if a state change should occur
   * Guard functions have access to the current `State` and the triggering `Event`
   * Guard functions return a `bool`, where `true` means the state change can occur
 * `IsFinal` will stop the machine if the state change is successful
 * `OnFail` is a function that is called when the `Guard` function returns `false` and no state change occurs
   * State change failure functions have access to the current `State` and the triggering `Event`
 * `OnSuccess` is a function that is called when the `Guard` function returns `true` and a state change occurred
   * State change success functions have access to the current `State` and the triggering `Event`
 * `To` is the next state the machine will enter if the state change is successful

Create a transition struct from the `Middle` state to the `End` state, triggered by the `WorkComplete` event.

```go
stt[Middle][WorkComplete] = &cism.Transition{
    Guard: func(s cism.State, e cism.Event) bool {
        return true
    },
    IsFinal: true,
    OnFail: func(s cism.State, e cism.Event) {},
    OnSuccess: func(s cism.State, e cism.Event) {},
    To: End,
}
```

### State Machine

The state machine is responsible for storing the current state and handling state change events.
It operates off of a state transition table to know what states to transition to for a sent event.
The machine can only be in one state at a given time.
It can be used to represent finite and infinite states.

Create a state machine struct with the state transition table.

```go
machine := &cism.Machine{
    States: stt,
}
```

 * `States` is a state transition table that the machine uses to determine how and when to transition

#### Start

Start the state machine with an initial state.

```go
err := machine.Start(Begin)
```

The returned error could be one four types of errors that will result in the machine not starting.
If the machine's `States` property is an empty `StateTransitionTable`, `Start` will return `ErrMissingStates`.
If the `State` is not defined in the `StateTransitionTable`, `Start` will return `ErrStateNotDefined`.
If the machine has been stopped, `Start` will return `ErrMachineStopped`.
If the machine has already been started, `Start` will return `ErrMachineStarted`.

A successful invocation of `Start` will flag the machine as started.
It will set the current state to the passed in initial `State`.

#### Send Event

Send an event to the machine to attempt a transition resulting in a state change.

```go
err := machine.Send(SetupDone)
```

The returned error could be one of three types of errors that will result in no attempted transition.
If the machine has not been started, `Send` will return `ErrMachineNotStarted`.
If the machine has been stopped, `Send` will return `ErrMachineStopped`.
If no transition is defined for the `Event` in the current state in the `StateTransitionTable`, `Send` will return
`ErrMissingTransition`.

When none of the above errors are encountered, `Send` will tell the machine to attempt a transition.
Refer to the `Transition` section for details on the lifecycle of a state change.
A successful invocation of `Send` will set the current state to the `Transition`'s `To` property's `State`.

#### Stop

Stop the machine, effectively preventing any new state changes.

```go
machine.Stop()
```

`Stop` flags the machine as stopped, and no new state changes can occur.
A machine can be stopped multiple times without error.

#### Reset

Reset the machine, allowing a stopped machine to be started again.

```go
err := machine.Reset()
```

The returned error will be `ErrMachineNotStopped` if the machine has already been stopped.

If the machine has been stopped, `Reset` will set the current state to the initial state set when `Start` was called.
The machine will be flagged as not started and not stopped.
It will clear the history log, as well.

#### Current State

Get the current state the machine is in.

```go
curr := machine.Current()
```

#### History Log

Get the history log of past states and their triggering events.

```go
hist := machine.History()
```

`History` will return a copy of the machine's `[]HistoryRecord` internal log.
The `HistoryRecord` struct is essentially a tuple of a `State` and `Event`.
Any modification to this history log copy will not affect the machine's actual history log it maintains.

## Example

The following example shows a simple state machine setup using `CISM`.
Define states and events to be used in the state transition table.
Create a state transition table for the state machine.
Create the state machine.
Then, send events and trigger state changes.
Finally, stop the machine and reset it.

```go
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
```

## License

_CISM_ is [MIT licensed](LICENSE).
