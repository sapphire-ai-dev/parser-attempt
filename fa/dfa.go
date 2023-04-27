package fa

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

// DFA deterministic finite automaton, state 0 is the start state
type DFA struct {
	InputNames []string
	States     []*DFAState
}

func (a *DFA) addTransition(src, input, dst int) {
	a.States[src].Transitions[input] = dst
}

func (a *DFA) Valid(state int) bool {
	return state >= 0 && state < len(a.States)
}

// Step returns [move success, new state is accepting, new state ID]
func (a *DFA) Step(src, input int) (bool, bool, int) {
	if a.Valid(src) {
		if dst, seen := a.States[src].Transitions[input]; seen {
			return true, a.States[dst].Accepts, dst
		}
	}

	return false, false, -1
}

func (a *DFA) Debug() string {
	result := "States\n"
	for src, state := range a.States {
		result += fmt.Sprintf("%s%d: %s\n", util.Tab, src, state.Name)
	}

	result += "Transitions:\n"
	for src, state := range a.States {
		for input, dst := range state.Transitions {
			result = result + fmt.Sprintf("%s%d => %s => %d\n", util.Tab, src, a.InputNames[input], dst)
		}
	}

	return result
}

// NewDFA transitions: list of [src, input, dst]
func NewDFA(inputNames, stateNames []string, accepts []bool, transitions [][3]int) *DFA {
	result := &DFA{InputNames: inputNames}
	for i := range stateNames {
		result.States = append(result.States, newDFAState(stateNames[i], accepts[i]))
	}

	for _, transition := range transitions {
		result.addTransition(transition[0], transition[1], transition[2])
	}

	return result
}

type DFAState struct {
	Name        string
	Accepts     bool
	Transitions map[int]int
}

func newDFAState(name string, accepts bool) *DFAState {
	return &DFAState{
		Name:        name,
		Accepts:     accepts,
		Transitions: map[int]int{},
	}
}
