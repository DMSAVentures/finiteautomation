package internal

import (
	"errors"
	"fmt"
)

var (
	ErrNoTransition      = errors.New("no transition defined")
	ErrInvalidInput      = errors.New("invalid input symbol")
	ErrNoStates          = errors.New("no states defined")
	ErrNoAlphabet        = errors.New("no alphabet defined")
	ErrInvalidInitial    = errors.New("initial state not in states list")
	ErrInvalidFinal      = errors.New("final state not in states list")
	ErrInvalidTransition = errors.New("transition references undefined state")
)

// FSM represents a generic finite state machine
type FSM[S State, I Input] struct {
	config        FSMConfig[S, I]
	currentState  S
	stateHistory  []S
	alphabetSet   map[I]struct{}
	finalStateSet map[S]struct{}
}

// validateAndBuildLookupSets validates the config and builds O(1) lookup sets
func validateAndBuildLookupSets[S State, I Input](config FSMConfig[S, I]) (
	alphabetSet map[I]struct{},
	finalStateSet map[S]struct{},
	err error,
) {
	// Validate basic requirements
	if len(config.States) == 0 {
		return nil, nil, ErrNoStates
	}
	if len(config.Alphabet) == 0 {
		return nil, nil, ErrNoAlphabet
	}

	// Build state set for validation (temporary, not returned)
	stateSet := make(map[S]struct{}, len(config.States))
	for _, state := range config.States {
		stateSet[state] = struct{}{}
	}

	// Validate initial state
	if _, exists := stateSet[config.InitialState]; !exists {
		return nil, nil, ErrInvalidInitial
	}

	// Validate final states (optional)
	if len(config.FinalStates) > 0 {
		for _, finalState := range config.FinalStates {
			if _, exists := stateSet[finalState]; !exists {
				return nil, nil, fmt.Errorf("%w: %v", ErrInvalidFinal, finalState)
			}
		}
	}

	// Build alphabet set for O(1) lookups
	alphabetSet = make(map[I]struct{}, len(config.Alphabet))
	for _, input := range config.Alphabet {
		alphabetSet[input] = struct{}{}
	}

	// Validate transitions
	for fromState, transitions := range config.Transitions {
		if _, exists := stateSet[fromState]; !exists {
			return nil, nil, fmt.Errorf("%w: from state %v", ErrInvalidTransition, fromState)
		}
		for input, toState := range transitions {
			if _, exists := alphabetSet[input]; !exists {
				return nil, nil, fmt.Errorf("%w: input %v not in alphabet", ErrInvalidInput, input)
			}
			if _, exists := stateSet[toState]; !exists {
				return nil, nil, fmt.Errorf("%w: to state %v", ErrInvalidTransition, toState)
			}
		}
	}

	// Build final state set for O(1) IsInFinalState lookups
	finalStateSet = make(map[S]struct{}, len(config.FinalStates))
	for _, state := range config.FinalStates {
		finalStateSet[state] = struct{}{}
	}

	return alphabetSet, finalStateSet, nil
}

// NewFSM creates a new FSM from a configuration
func NewFSM[S State, I Input](config FSMConfig[S, I]) (*FSM[S, I], error) {
	alphabetSet, finalStateSet, err := validateAndBuildLookupSets(config)
	if err != nil {
		return nil, err
	}

	return &FSM[S, I]{
		config:        config,
		currentState:  config.InitialState,
		stateHistory:  []S{config.InitialState},
		alphabetSet:   alphabetSet,
		finalStateSet: finalStateSet,
	}, nil
}

// Reset resets the FSM to its initial state
func (f *FSM[S, I]) Reset() {
	f.currentState = f.config.InitialState
	f.stateHistory = []S{f.currentState}
}

// CurrentState returns the current state
func (f *FSM[S, I]) CurrentState() S {
	return f.currentState
}

// StateHistory returns the state transition history
func (f *FSM[S, I]) StateHistory() []S {
	history := make([]S, len(f.stateHistory))
	copy(history, f.stateHistory)
	return history
}

// Transition processes a single input and transitions to the next state
func (f *FSM[S, I]) Transition(input I) error {
	stateTransitions, exists := f.config.Transitions[f.currentState]
	if !exists {
		return fmt.Errorf("%w for state: %v", ErrNoTransition, f.currentState)
	}

	nextState, exists := stateTransitions[input]
	if !exists {
		return fmt.Errorf("%w from state %v with input %v", ErrNoTransition, f.currentState, input)
	}

	f.currentState = nextState
	f.stateHistory = append(f.stateHistory, nextState)
	return nil
}

// Execute processes a sequence of inputs, resetting the FSM first
func (f *FSM[S, I]) Execute(inputs []I) error {
	f.Reset()
	return f.Process(inputs)
}

// Process processes a sequence of inputs without resetting, continuing from current state
func (f *FSM[S, I]) Process(inputs []I) error {
	for i, input := range inputs {
		if err := f.Transition(input); err != nil {
			return fmt.Errorf("error at position %d: %w", i, err)
		}
	}
	return nil
}

// IsInFinalState returns true if the current state is a final state
func (f *FSM[S, I]) IsInFinalState() bool {
	if len(f.finalStateSet) == 0 {
		return false
	}
	_, exists := f.finalStateSet[f.currentState]
	return exists
}

// ValidateInput checks if an input is in the alphabet
func (f *FSM[S, I]) ValidateInput(input I) bool {
	_, exists := f.alphabetSet[input]
	return exists
}
