package finiteautomation

type State comparable
type Input comparable

// FSMConfig defines the configuration for a generic FSM
type FSMConfig[S State, I Input] struct {
	States       []S
	Alphabet     []I
	InitialState S
	FinalStates  []S
	Transitions  map[S]map[I]S
}

// AddTransition adds a transition to the FSM configuration
func (cfg *FSMConfig[S, I]) AddTransition(from S, input I, to S) {
	if cfg.Transitions == nil {
		cfg.Transitions = make(map[S]map[I]S)
	}
	if cfg.Transitions[from] == nil {
		cfg.Transitions[from] = make(map[I]S)
	}
	cfg.Transitions[from][input] = to
}
