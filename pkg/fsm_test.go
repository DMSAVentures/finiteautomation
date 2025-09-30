package internal

import (
	"testing"
)

// Test with string states and rune inputs
func TestGenericFSM_StringStates(t *testing.T) {
	config := FSMConfig[string, rune]{
		States:       []string{"locked", "unlocked"},
		Alphabet:     []rune{'C', 'P'}, // Coin, Push
		InitialState: "locked",
		FinalStates:  []string{"unlocked"},
		Transitions: map[string]map[rune]string{
			"locked": {
				'C': "unlocked",
				'P': "locked",
			},
			"unlocked": {
				'C': "unlocked",
				'P': "locked",
			},
		},
	}

	fsm, _ := NewFSM(config)

	tests := []struct {
		name          string
		input         rune
		expectedState string
	}{
		{"initial state", 0, "locked"}, // 0 means no transition
		{"coin inserted", 'C', "unlocked"},
		{"push from unlocked", 'P', "locked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != 0 {
				fsm.Transition(tt.input)
			}
			if fsm.CurrentState() != tt.expectedState {
				t.Errorf("Expected state %s, got %s", tt.expectedState, fsm.CurrentState())
			}
		})
	}
}

// Test with integer states and string inputs
func TestGenericFSM_IntStates(t *testing.T) {
	config := FSMConfig[int, string]{
		States:       []int{0, 1, 2},
		Alphabet:     []string{"next", "prev", "reset"},
		InitialState: 0,
		FinalStates:  []int{2},
		Transitions: map[int]map[string]int{
			0: {"next": 1, "prev": 0, "reset": 0},
			1: {"next": 2, "prev": 0, "reset": 0},
			2: {"next": 2, "prev": 1, "reset": 0},
		},
	}

	fsm, _ := NewFSM(config)

	tests := []struct {
		inputs   []string
		expected int
		isFinal  bool
	}{
		{[]string{}, 0, false},
		{[]string{"next"}, 1, false},
		{[]string{"next", "next"}, 2, true},
		{[]string{"next", "next", "prev"}, 1, false},
		{[]string{"next", "next", "reset"}, 0, false},
	}

	for _, tt := range tests {
		fsm.Execute(tt.inputs)
		if fsm.CurrentState() != tt.expected {
			t.Errorf("After %v, expected state %d, got %d", tt.inputs, tt.expected, fsm.CurrentState())
		}
		if fsm.IsInFinalState() != tt.isFinal {
			t.Errorf("After %v, IsFinal should be %v, got %v", tt.inputs, tt.isFinal, fsm.IsInFinalState())
		}
	}
}

// Test with custom struct types
type TestState struct {
	Name  string
	Value int
}

type TestInput struct {
	Action string
}

func TestGenericFSM_StructTypes(t *testing.T) {
	stateA := TestState{"A", 1}
	stateB := TestState{"B", 2}

	inputX := TestInput{"X"}
	inputY := TestInput{"Y"}

	config := FSMConfig[TestState, TestInput]{
		States:       []TestState{stateA, stateB},
		Alphabet:     []TestInput{inputX, inputY},
		InitialState: stateA,
		FinalStates:  []TestState{stateB},
		Transitions: map[TestState]map[TestInput]TestState{
			stateA: {
				inputX: stateB,
				inputY: stateA,
			},
			stateB: {
				inputX: stateA,
				inputY: stateB,
			},
		},
	}

	fsm, _ := NewFSM(config)

	tests := []struct {
		name          string
		input         *TestInput
		expectedState TestState
		isFinal       bool
	}{
		{"initial state", nil, stateA, false},
		{"transition X to B", &inputX, stateB, true},
		{"transition X back to A", &inputX, stateA, false},
		{"transition Y stays at A", &inputY, stateA, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != nil {
				fsm.Transition(*tt.input)
			}
			if fsm.CurrentState() != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, fsm.CurrentState())
			}
			if fsm.IsInFinalState() != tt.isFinal {
				t.Errorf("Expected IsFinal=%v, got %v", tt.isFinal, fsm.IsInFinalState())
			}
		})
	}
}

func TestGenericFSM_ValidateInput(t *testing.T) {
	config := FSMConfig[string, int]{
		States:       []string{"A"},
		Alphabet:     []int{1, 2, 3},
		InitialState: "A",
		FinalStates:  []string{},
		Transitions:  map[string]map[int]string{},
	}

	fsm, _ := NewFSM(config)

	tests := []struct {
		name  string
		input int
		valid bool
	}{
		{"valid input 1", 1, true},
		{"valid input 2", 2, true},
		{"valid input 3", 3, true},
		{"invalid input 4", 4, false},
		{"invalid input 0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if fsm.ValidateInput(tt.input) != tt.valid {
				t.Errorf("ValidateInput(%d) = %v, want %v", tt.input, !tt.valid, tt.valid)
			}
		})
	}
}

func TestFSMValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    FSMConfig[string, rune]
		wantError bool
	}{
		{
			name: "valid config",
			config: FSMConfig[string, rune]{
				States:       []string{"A", "B"},
				Alphabet:     []rune{'0', '1'},
				InitialState: "A",
				FinalStates:  []string{"B"},
				Transitions: map[string]map[rune]string{
					"A": {'0': "B"},
				},
			},
			wantError: false,
		},
		{
			name: "no states",
			config: FSMConfig[string, rune]{
				States:   []string{},
				Alphabet: []rune{'0'},
			},
			wantError: true,
		},
		{
			name: "no alphabet",
			config: FSMConfig[string, rune]{
				States:   []string{"A"},
				Alphabet: []rune{},
			},
			wantError: true,
		},
		{
			name: "invalid initial state",
			config: FSMConfig[string, rune]{
				States:       []string{"A", "B"},
				Alphabet:     []rune{'0'},
				InitialState: "C",
			},
			wantError: true,
		},
		{
			name: "invalid final state",
			config: FSMConfig[string, rune]{
				States:       []string{"A", "B"},
				Alphabet:     []rune{'0'},
				InitialState: "A",
				FinalStates:  []string{"C"},
			},
			wantError: true,
		},
		{
			name: "transition to undefined state",
			config: FSMConfig[string, rune]{
				States:       []string{"A"},
				Alphabet:     []rune{'0'},
				InitialState: "A",
				Transitions: map[string]map[rune]string{
					"A": {'0': "B"}, // B not in states
				},
			},
			wantError: true,
		},
		{
			name: "transition with undefined input",
			config: FSMConfig[string, rune]{
				States:       []string{"A", "B"},
				Alphabet:     []rune{'0'},
				InitialState: "A",
				Transitions: map[string]map[rune]string{
					"A": {'1': "B"}, // '1' not in alphabet
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFSM(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("NewFSM() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestGenericFSM_StateHistory(t *testing.T) {
	config := FSMConfig[string, rune]{
		States:       []string{"A", "B", "C"},
		Alphabet:     []rune{'0', '1'},
		InitialState: "A",
		FinalStates:  []string{"C"},
		Transitions: map[string]map[rune]string{
			"A": {'0': "B", '1': "C"},
			"B": {'0': "C", '1': "A"},
			"C": {'0': "A", '1': "B"},
		},
	}

	tests := []struct {
		name     string
		inputs   []rune
		expected []string
	}{
		{"no transitions", []rune{}, []string{"A"}},
		{"single transition", []rune{'0'}, []string{"A", "B"}},
		{"multiple transitions", []rune{'0', '1', '0'}, []string{"A", "B", "A", "B"}},
		{"longer sequence", []rune{'1', '1', '0'}, []string{"A", "C", "B", "C"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm, _ := NewFSM(config)
			fsm.Execute(tt.inputs)
			history := fsm.StateHistory()

			if len(history) != len(tt.expected) {
				t.Fatalf("History length = %d, want %d", len(history), len(tt.expected))
			}

			for i, state := range tt.expected {
				if history[i] != state {
					t.Errorf("History[%d] = %s, want %s", i, history[i], state)
				}
			}
		})
	}
}
