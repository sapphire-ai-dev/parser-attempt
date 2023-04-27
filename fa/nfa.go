package fa

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/util"
	"sort"
)

// NFA non-deterministic finite automaton, state 0 is the start state
type NFA struct {
	InputNames []string
	States     []*NFAState
}

func (a *NFA) addTransition(src, input, dst int) {
	if src < 0 || src >= len(a.States) || dst < 0 || dst >= len(a.States) {
		panic("invalid src or dst")
	}

	if input >= 0 {
		if _, seen := a.States[src].Transitions[input]; !seen {
			a.States[src].Transitions[input] = map[int]bool{}
		}

		a.States[src].Transitions[input][dst] = true
	} else { // negative transition input implies epsilon transition
		a.States[src].EpsTransitions[dst] = true
	}
}

func (a *NFA) finalizeTransitions() {
	// currently O(|state|^2), can probably improve to O(|state|)
	for i, state := range a.States {
		epsTransitions := map[int]bool{}
		a.finalizeEpsTransitionsHelper(i, epsTransitions)
		state.EpsTransitions = epsTransitions
	}

	// currently O(|state|^4), not sure if possible to improve
	for src := range a.States {
		for input := range a.States[src].Transitions {
			// going through input from src can be expanded as three steps:
			// 1. going through all epsilon transitions from src
			// 2. going through input transition
			// 3. going through all epsilon transitions again
			epsTransitions := map[int]bool{}
			for eps := range a.States[src].EpsTransitions { // step 1
				epsTransitions = util.KeyUnion(epsTransitions, a.States[eps].Transitions[input]) // step 2
			}

			// collect at step 2 to reduce the whole process from O(|state|^3) to O(|state|^2)
			newTransitions := util.MapCopy[int, bool](epsTransitions)
			for transition := range epsTransitions {
				newTransitions = util.KeyUnion(newTransitions, a.States[transition].EpsTransitions) // step 3
			}
			a.States[src].Transitions[input] = newTransitions
		}
	}
}

func (a *NFA) finalizeEpsTransitionsHelper(curr int, done map[int]bool) {
	if done[curr] {
		return
	}

	done[curr] = true
	for eps := range a.States[curr].EpsTransitions {
		a.finalizeEpsTransitionsHelper(eps, done)
	}
}

func (a *NFA) Valid(states map[int]bool) bool {
	for state := range states {
		if state < 0 || state >= len(a.States) {
			return false
		}
	}

	return true
}

// Step returns [move success, new state is accepting, set of new state IDs]
func (a *NFA) Step(srcs map[int]bool, input int) (bool, bool, map[int]bool) {
	if !a.Valid(srcs) {
		return false, false, map[int]bool{}
	}

	dsts := map[int]bool{}
	for src := range srcs {
		dsts = util.KeyUnion(dsts, a.States[src].Transitions[input])
	}

	accepted := false
	for dst := range dsts {
		if a.States[dst].Accepts {
			accepted = true
			break
		}
	}

	return len(dsts) > 0, accepted, dsts
}

// NewNFA transitions: list of [src, input, dst]
func NewNFA(inputNames, stateNames []string, accepts []bool, transitions [][3]int) *NFA {
	result := &NFA{InputNames: inputNames}
	for i := range stateNames {
		result.States = append(result.States, newNFAState(stateNames[i], accepts[i]))
	}

	for _, transition := range transitions {
		result.addTransition(transition[0], transition[1], transition[2])
	}

	result.finalizeTransitions()
	return result
}

type NFAState struct {
	Name           string
	Accepts        bool
	Transitions    map[int]map[int]bool // input ID -> set of next state IDs
	EpsTransitions map[int]bool         // next states
}

func newNFAState(name string, accepts bool) *NFAState {
	return &NFAState{
		Name:           name,
		Accepts:        accepts,
		Transitions:    map[int]map[int]bool{},
		EpsTransitions: map[int]bool{},
	}
}

// CompileDFA returns the compiled DFA and a state ID mapping: DFA state ID -> set of NFA state IDs
func (a *NFA) CompileDFA() (*DFA, []map[int]bool) {
	dfaStateNames := map[string]int{}
	acceptStates := map[int]bool{}
	transitions := map[int][3]int{}
	stateIdMapping := map[int]map[int]bool{} // DFA state ID -> set of NFA state IDs
	a.compileDFAHelper(dfaStateNames, acceptStates, transitions, stateIdMapping, a.States[0].EpsTransitions)

	dfaStateList := make([]string, len(dfaStateNames))
	acceptList := make([]bool, len(dfaStateNames))
	for stateName, stateId := range dfaStateNames {
		dfaStateList[stateId] = stateName
		acceptList[stateId] = acceptStates[stateId]
	}

	transitionList := make([][3]int, len(transitions))
	for i, transition := range transitions {
		transitionList[i] = transition
	}

	stateIdMappingList := make([]map[int]bool, len(dfaStateNames))
	for dfaState, nfaStates := range stateIdMapping {
		stateIdMappingList[dfaState] = nfaStates
	}

	return NewDFA(a.InputNames, dfaStateList, acceptList, transitionList), stateIdMappingList
}

// transitions made into a map instead of a list to pass by reference
func (a *NFA) compileDFAHelper(stateNames map[string]int, acceptStates map[int]bool,
	transitions map[int][3]int, stateIdMapping map[int]map[int]bool, currStates map[int]bool) string {
	currStatesName := a.compileDFAGenName(currStates)
	if _, seen := stateNames[currStatesName]; seen {
		return currStatesName
	}

	stateNames[currStatesName] = len(stateNames)
	stateIdMapping[stateNames[currStatesName]] = currStates
	for input := range a.InputNames {
		valid, accepts, dsts := a.Step(currStates, input)
		if !valid {
			continue
		}

		dstsName := a.compileDFAHelper(stateNames, acceptStates, transitions, stateIdMapping, dsts)
		transitions[len(transitions)] = [3]int{stateNames[currStatesName], input, stateNames[dstsName]}
		if accepts {
			acceptStates[stateNames[dstsName]] = true
		}
	}

	return currStatesName
}

func (a *NFA) compileDFAGenName(currStates map[int]bool) string {
	stateList := util.KeyList(currStates)
	sort.Ints(stateList)
	var stateNames []string
	for _, state := range stateList {
		stateNames = append(stateNames, a.States[state].Name)
	}

	return fmt.Sprintf("%v", stateNames)
}
