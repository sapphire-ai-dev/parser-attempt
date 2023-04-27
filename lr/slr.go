package lr

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/cfg"
	"github.com/sapphire-ai-dev/parser-attempt/fa"
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

type SLRParser struct {
	// permanent fields
	cfg       *cfg.CFG
	prefixDFA *fa.DFA
	actions   []map[int]int // prefixDFA state ID -> { input ID -> action ID }
	reduces   []map[int]*reduceMove

	// temporary fields to track parsing progress
	parseInput []int
	parsePos   int
	parseStack []*parseState
	parseNodes []*cfg.ParseTreeNode
}

func (p *SLRParser) Debug() string {
	inputNames := util.InvertMap(p.cfg.Symbols.Names)
	inputNames[-2] = "$"
	result := "PrefixDFA\n"
	result += p.prefixDFA.Debug()
	result += "Actions\n"
	for state := range p.actions {
		for input, action := range p.actions[state] {
			actionStr := "shift"
			if action == ActionIdReduce {
				actionStr = fmt.Sprintf("reduce %v => %s", util.MapOverList(inputNames,
					p.reduces[state][input].from), inputNames[p.reduces[state][input].to])
			}

			result += fmt.Sprintf("%s%d %s: %s\n", util.Tab, state, inputNames[input], actionStr)
		}
	}

	return result
}

func newSLRParser(g *cfg.CFG) *SLRParser {
	return &SLRParser{cfg: g}
}

type reduceMove struct {
	from []int
	to   int
}

func (m *reduceMove) match(o *reduceMove) bool {
	if len(m.from) != len(o.from) {
		return false
	}

	for i := range m.from {
		if m.from[i] != o.from[i] {
			return false
		}
	}

	if m.to != o.to {
		return false
	}

	return true
}

func newReduceMove(from []int, to int) *reduceMove {
	return &reduceMove{
		from: from,
		to:   to,
	}
}
