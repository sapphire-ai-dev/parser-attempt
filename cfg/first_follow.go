package cfg

import (
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

func (g *CFG) buildFirstSet() {
	var firstSet []map[int]bool
	for range g.Symbols.Names {
		firstSet = append(firstSet, map[int]bool{})
	}

	for terminal := range g.Symbols.Terminals {
		firstSet[terminal][terminal] = true
	}

	invNames := util.InvertMap(g.Symbols.Names)
	invNames[-1] = "eps"

	// some weird variation of bellman ford
	addedNew := true
	for addedNew {
		addedNew = false
		for nonTerminal, rules := range g.Rules {
			for _, rule := range rules {
				bypassed := 0

				for i := 0; i < len(rule); i++ {
					setSize := len(firstSet[nonTerminal])
					// can optimize this to record change in each first set, skip this step if no change recorded
					firstSet[nonTerminal] = util.KeyUnion(firstSet[nonTerminal], firstSet[rule[i]])
					addedNew = addedNew || setSize < len(firstSet[nonTerminal])

					if firstSet[rule[i]][-1] {
						bypassed++
					} else {
						break
					}
				}

				if bypassed == len(rule) {
					addedNew = util.AddedNew(firstSet[nonTerminal], -1, true) || addedNew
				}
			}
		}
	}

	g.FirstSet = firstSet
}

func (g *CFG) buildFollowSet() {
	var followSet []map[int]bool
	for range g.Symbols.Names {
		followSet = append(followSet, map[int]bool{})
	}

	followSet[0][-2] = true // add $ to follow(S)

	// a, b are sequences of symbols, A, B are symbols
	// for rule B -> aAb, add first(b) to follow(A)
	for _, rules := range g.Rules {
		for _, rule := range rules {
			for i := 1; i < len(rule); i++ {
				for elem := range g.FirstSet[rule[i]] {
					if elem != -1 {
						followSet[rule[i-1]][elem] = true
					}
				}
			}
		}
	}

	// for rule B -> aAb, if |b| == 0 or first(b) contains eps, add follow(B) to follow(A)
	addedNew := true
	for addedNew {
		addedNew = false
		for nonTerminal, rules := range g.Rules {
			for _, rule := range rules {
				for i := len(rule) - 1; i >= 0; i-- {
					setSize := len(followSet[rule[i]])
					followSet[rule[i]] = util.KeyUnion(followSet[rule[i]], followSet[nonTerminal])
					addedNew = addedNew || setSize < len(followSet[rule[i]])

					if !g.FirstSet[rule[i]][-1] {
						break
					}
				}
			}
		}
	}

	g.FollowSet = followSet
}
