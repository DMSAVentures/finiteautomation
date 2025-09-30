package main

import (
	"finiteautomation/internal"
	"fmt"
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
	fsm *internal.FSM[ModState, BinarySymbol]
}

// NewModThreeGeneric creates a new modulo-3 FSM using generics
func NewModThreeGeneric() *ModThreeGeneric {
	config := internal.FSMConfig[ModState, BinarySymbol]{
		States:       []ModState{ModState0, ModState1, ModState2},
		Alphabet:     []BinarySymbol{Binary0, Binary1},
		InitialState: ModState0,
		FinalStates:  []ModState{ModState0},
		Transitions: map[ModState]map[BinarySymbol]ModState{
			ModState0: {
				Binary0: ModState0, // (0*2 + 0) mod 3 = 0
				Binary1: ModState1, // (0*2 + 1) mod 3 = 1
			},
			ModState1: {
				Binary0: ModState2, // (1*2 + 0) mod 3 = 2
				Binary1: ModState0, // (1*2 + 1) mod 3 = 0
			},
			ModState2: {
				Binary0: ModState1, // (2*2 + 0) mod 3 = 1
				Binary1: ModState2, // (2*2 + 1) mod 3 = 2
			},
		},
	}

	fsm, err := internal.NewFSM(config)
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

// IsDivisibleByThree checks if a binary number is divisible by 3
func (m *ModThreeGeneric) IsDivisibleByThree(binaryStr string) bool {
	inputs, err := m.ParseInput(binaryStr)
	if err != nil {
		return false
	}
	if inputs == nil {
		// Empty string or "0" is considered divisible by 3
		return true
	}

	err = m.fsm.Execute(inputs)
	if err != nil {
		return false
	}

	return true
}
