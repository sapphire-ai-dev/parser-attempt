package lr

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/cfg"
	"github.com/sapphire-ai-dev/parser-attempt/fa"
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

func BuildSLRParser(g *cfg.CFG) *SLRParser {
	parser := newSLRParser(g)
	var stateMapping []map[int]bool
	prefixNFA := parser.buildPrefixNFA()
	parser.prefixDFA, stateMapping = prefixNFA.CompileDFA()
	parser.buildActions(stateMapping)
	fmt.Println(parser.prefixDFA.Debug())
	return parser
}

func (p *SLRParser) buildPrefixNFA() *fa.NFA {
	invNames := util.InvertMap(p.cfg.Symbols.Names)
	nfaStateNames := map[string]int{}
	ruleStartStates := map[int][]int{} // nonTerminalId -> list of rule start nfa states
	for nonTerminal := range p.cfg.Rules {
		for _, rule := range p.cfg.Rules[nonTerminal] {
			names := util.MapOverList(invNames, rule)
			ruleStartStates[nonTerminal] = append(ruleStartStates[nonTerminal], len(nfaStateNames))
			for i := 0; i <= len(rule); i++ {
				stateName := fmt.Sprintf("%s => %v%v", invNames[nonTerminal], names[:i], names[i:])
				nfaStateNames[stateName] = len(nfaStateNames)
			}
		}
	}

	var transitions [][3]int
	nfaState := 0 // used to retrace nfa states the same way they were created
	for nonTerminal := range p.cfg.Rules {
		for _, rule := range p.cfg.Rules[nonTerminal] {
			for i := 0; i <= len(rule); i++ {
				if i < len(rule) {
					transitions = append(transitions, [3]int{nfaState, rule[i], nfaState + 1}) // shift
					for _, reduceStateId := range ruleStartStates[rule[i]] {
						transitions = append(transitions, [3]int{nfaState, -1, reduceStateId})
					}
				}
				nfaState++
			}
		}
	}

	var accepts []bool
	for range nfaStateNames {
		accepts = append(accepts, true)
	}

	inputNameList := make([]string, len(invNames))
	stateNameList := make([]string, len(nfaStateNames))
	for iId, inputName := range invNames {
		inputNameList[iId] = inputName
	}
	for stateName, sId := range nfaStateNames {
		stateNameList[sId] = stateName
	}

	return fa.NewNFA(inputNameList, stateNameList, accepts, transitions)
}

func (p *SLRParser) buildActions(stateIdMapping []map[int]bool) {
	for range p.prefixDFA.States {
		p.actions = append(p.actions, map[int]int{})
		p.reduces = append(p.reduces, map[int]*reduceMove{})
	}

	invStateMapping := map[int]map[int]bool{}
	for dfaState, nfaStates := range stateIdMapping {
		for nfaState := range nfaStates {
			if _, seen := invStateMapping[nfaState]; !seen {
				invStateMapping[nfaState] = map[int]bool{}
			}

			invStateMapping[nfaState][dfaState] = true
		}
	}

	nfaState := 0 // used to retrace nfa states the same way they were created
	for nonTerminal := range p.cfg.Rules {
		for _, rule := range p.cfg.Rules[nonTerminal] {
			// shifts
			for _, input := range rule {
				for dfaSrc := range invStateMapping[nfaState] {
					if valid, _, _ := p.prefixDFA.Step(dfaSrc, input); valid {
						p.addShiftAction(dfaSrc, input)
					}
				}
				nfaState++
			}

			// reduces
			for dfaSrc := range invStateMapping[nfaState] {
				for input := range p.prefixDFA.InputNames {
					if p.cfg.FollowSet[nonTerminal][input] {
						p.addReduceAction(dfaSrc, input, rule, nonTerminal)
					}
				}

				if p.cfg.FollowSet[nonTerminal][-2] {
					p.addReduceAction(dfaSrc, -2, rule, nonTerminal)
				}
			}
			nfaState++
		}
	}
}

func (p *SLRParser) addShiftAction(dfaSrc, input int) {
	if actionId, seen := p.actions[dfaSrc][input]; seen && actionId == ActionIdReduce {
		panic("shift-reduce collision")
	}

	p.actions[dfaSrc][input] = ActionIdShift
}

func (p *SLRParser) addReduceAction(dfaSrc, input int, from []int, to int) {
	move := newReduceMove(from, to)
	if actionId, seen := p.actions[dfaSrc][input]; seen {
		if actionId == ActionIdShift {
			panic("reduce-shift collision")
		}

		if !p.reduces[dfaSrc][input].match(move) {
			panic("reduce-reduce collision")
		}
	}

	p.actions[dfaSrc][input] = ActionIdReduce
	p.reduces[dfaSrc][input] = move
}
