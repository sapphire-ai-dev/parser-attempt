package lr

import (
	"github.com/sapphire-ai-dev/parser-attempt/cfg"
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

func (p *SLRParser) ParseRaw(input []string) *cfg.ParseTreeNode {
	return p.Parse(util.MapOverList(p.cfg.Symbols.Names, input))
}

func (p *SLRParser) Parse(input []int) *cfg.ParseTreeNode {
	input = append(input, -2) // add $ at end of input
	p.parseStateReset(input)
	for len(p.parseStack) > 1 || input[p.parsePos] != -2 {
		if !p.action() {
			panic("invalid action")
		}
	}

	if len(p.parseNodes) != 1 {
		return nil
	}

	return p.parseNodes[0]
}

func (p *SLRParser) parseStateReset(input []int) {
	p.parseInput = input
	p.parseStack = []*parseState{newParseState(-1, 0)}
	p.parseNodes = []*cfg.ParseTreeNode{}
	p.parsePos = 0
}

const (
	ActionIdShift = iota
	ActionIdReduce
)

func (p *SLRParser) action() bool {
	currState := p.parseStack[len(p.parseStack)-1].state
	nextSymbol := p.parseInput[p.parsePos]
	if actionId, seen := p.actions[currState][nextSymbol]; seen {
		if actionId == ActionIdShift {
			p.shift(currState, nextSymbol)
			return true
		} else if actionId == ActionIdReduce {
			p.reduce(currState, p.reduces[currState][nextSymbol])
			return true
		}
	}

	return false
}

func (p *SLRParser) shift(currState, nextSymbol int) {
	p.parseNodes = append(p.parseNodes, cfg.NewParseTreeNode(nextSymbol))
	p.parsePos++

	valid, _, nextState := p.prefixDFA.Step(currState, nextSymbol)
	if !valid {
		panic("invalid shift")
	}

	p.parseStack = append(p.parseStack, newParseState(nextSymbol, nextState))
}

func (p *SLRParser) reduce(currState int, move *reduceMove) {
	newParseNode := cfg.NewParseTreeNode(move.to)
	cutoff := len(p.parseNodes) - len(move.from)
	for _, child := range p.parseNodes[cutoff:] {
		newParseNode.Children = append(newParseNode.Children, child)
	}
	p.parseNodes = p.parseNodes[:cutoff]
	p.parseNodes = append(p.parseNodes, newParseNode)

	p.parseStack = p.parseStack[:len(p.parseStack)-len(move.from)]
	currState = p.parseStack[len(p.parseStack)-1].state
	if currState == 0 && move.to == 0 {
		return
	}

	valid, _, nextState := p.prefixDFA.Step(currState, move.to)
	if !valid {
		panic("invalid reduce")
	}

	p.parseStack = append(p.parseStack, newParseState(move.to, nextState))
}

type parseState struct {
	symbol int
	state  int
}

func newParseState(symbol, state int) *parseState {
	return &parseState{
		symbol: symbol,
		state:  state,
	}
}
