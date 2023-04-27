package cfg

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

type SymbolSet struct {
	Terminals    map[int]bool   // set of terminal symbols (tokens)
	NonTerminals map[int]bool   // set of non-terminal symbols (expressions)
	Names        map[string]int // names used for referencing -> mapped integer ids
}

func (s *SymbolSet) valid() bool {
	return len(util.KeyIntersection(s.Terminals, s.NonTerminals)) == 0 &&
		util.KeyEqual(util.MapValSet(s.Names), util.KeyUnion(s.Terminals, s.NonTerminals))
}

func (s *SymbolSet) addSymbol(name string, isTerminal bool) int {
	if _, seen := s.Names[name]; seen {
		return s.Names[name]
	}

	symbolId := len(s.Names)
	if isTerminal {
		s.Terminals[symbolId] = true
	} else {
		s.NonTerminals[symbolId] = true
	}
	s.Names[name] = symbolId
	return symbolId
}

func (s *SymbolSet) debug() string {
	invNames := util.InvertMap(s.Names)
	var terminals, nonTerminals []string
	for symbol := range s.Terminals {
		terminals = append(terminals, invNames[symbol])
	}
	for symbol := range s.NonTerminals {
		nonTerminals = append(nonTerminals, invNames[symbol])
	}
	return fmt.Sprintf("%sTerminals: %v\n%sNon-Terminals: %v\n", util.Tab, terminals, util.Tab, nonTerminals)
}

type CFG struct {
	Symbols   *SymbolSet
	Rules     [][][]int
	FirstSet  []map[int]bool
	FollowSet []map[int]bool
}

func (g *CFG) addRule(lhs int, rhs []int) {
	g.Rules[lhs] = append(g.Rules[lhs], rhs)
}

func (g *CFG) Debug() string {
	invNames := util.InvertMap(g.Symbols.Names)
	invNames[-1] = "eps"
	invNames[-2] = "$"
	result := "Symbols:\n" + g.Symbols.debug()
	result += "Rules:\n"
	for symbol := range g.Rules {
		for _, rule := range g.Rules[symbol] {
			ruleStr := fmt.Sprintf("%s => %v", invNames[symbol], util.MapOverList(invNames, rule))
			result += fmt.Sprintf("%s%s\n", util.Tab, ruleStr)
		}
	}

	result += "First Sets:\n"
	for symbol, set := range g.FirstSet {
		result += fmt.Sprintln(util.Tab+invNames[symbol], util.MapOverList(invNames, util.KeyList(set)))
	}

	result += "Follow Sets:\n"
	for symbol, set := range g.FollowSet {
		result += fmt.Sprintln(util.Tab+invNames[symbol], util.MapOverList(invNames, util.KeyList(set)))
	}

	return result
}

func newCFG() *CFG {
	return &CFG{
		Symbols: &SymbolSet{
			Terminals:    map[int]bool{},
			NonTerminals: map[int]bool{},
			Names:        map[string]int{},
		},
	}
}
