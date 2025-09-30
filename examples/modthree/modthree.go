package modthree

import (
	"fmt"

	"github.com/dmsaventures/finiteautomation"
)

// ModState represents the states for modulo-3 computation
type ModState int

const (
	ModState0 ModState = 0 // Remainder 0
	ModState1 ModState = 1 // Remainder 1
	ModState2 ModState = 2 // Remainder 2
)

// BinarySymbol represents binary input symbols
type BinarySymbol rune

const (
	Binary0 BinarySymbol = '0'
	Binary1 BinarySymbol = '1'
)

// ModThreeGeneric uses the generic FSM for modulo-3 computation
type ModThreeGeneric struct {
	fsm *finiteautomation.FSM[ModState, BinarySymbol]
}

// NewModThreeGeneric creates a new modulo-3 FSM using generics
func NewModThreeGeneric() *ModThreeGeneric {
	// Based on the formal definition:
	// Q = (S0, S1, S2)
	// Σ = (0, 1)
	// q0 = S0
	// F = (S0, S1, S2)
	// δ(S0,0) = S0; δ(S0,1) = S1; δ(S1,0) = S2; δ(S1,1) = S0; δ(S2,0) = S1; δ(S2,1) = S2
	config := finiteautomation.FSMConfig[ModState, BinarySymbol]{
		States:       []ModState{ModState0, ModState1, ModState2}, // Q
		Alphabet:     []BinarySymbol{Binary0, Binary1},             // Σ
		InitialState: ModState0,                                    // q0
		FinalStates:  []ModState{ModState0, ModState1, ModState2},  // F - all states are accepting
		Transitions: map[ModState]map[BinarySymbol]ModState{       // δ
			ModState0: {
				Binary0: ModState0, // δ(S0,0) = S0
				Binary1: ModState1, // δ(S0,1) = S1
			},
			ModState1: {
				Binary0: ModState2, // δ(S1,0) = S2
				Binary1: ModState0, // δ(S1,1) = S0
			},
			ModState2: {
				Binary0: ModState1, // δ(S2,0) = S1
				Binary1: ModState2, // δ(S2,1) = S2
			},
		},
	}

	fsm, err := finiteautomation.NewFSM(config)
	if err != nil {
		// This shouldn't happen with our hardcoded config, but handle it anyway
		panic(fmt.Sprintf("Invalid ModThree FSM config: %v", err))
	}

	return &ModThreeGeneric{
		fsm: fsm,
	}
}

func (m *ModThreeGeneric) ParseInput(binaryStr string) ([]BinarySymbol, error) {
	if binaryStr == "" || binaryStr == "0" {
		return nil, nil
	}

	// Convert string to BinarySymbols
	inputs := make([]BinarySymbol, len(binaryStr))
	for i, char := range binaryStr {
		symbol := BinarySymbol(char)
		if !m.fsm.ValidateInput(symbol) {
			return nil, fmt.Errorf("invalid binary character '%c' at position %d", char, i)
		}
		inputs[i] = symbol
	}
	return inputs, nil
}

// ComputeModThree computes the modulo-3 remainder of a binary number
// Returns the remainder (0, 1, or 2) and whether the computation was successful
func (m *ModThreeGeneric) ComputeModThree(binaryStr string) (int, bool) {
	inputs, err := m.ParseInput(binaryStr)
	if err != nil {
		return 0, false
	}
	if inputs == nil {
		// Empty string or "0" has remainder 0
		return 0, true
	}

	err = m.fsm.Execute(inputs)
	if err != nil {
		return 0, false
	}

	// Return the current state which represents the remainder
	return int(m.fsm.CurrentState()), true
}

// IsDivisibleByThree checks if a binary number is divisible by 3
// A number is divisible by 3 if its remainder is 0 (i.e., final state is S0)
func (m *ModThreeGeneric) IsDivisibleByThree(binaryStr string) bool {
	remainder, ok := m.ComputeModThree(binaryStr)
	return ok && remainder == 0
}
