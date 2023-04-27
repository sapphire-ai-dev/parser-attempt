package cfg

import (
	"github.com/sapphire-ai-dev/parser-attempt/util"
)

type ParseTreeNode struct {
	Symbol   int
	Children []*ParseTreeNode
}

func (n *ParseTreeNode) Debug(g *CFG) string {
	names := util.InvertMap(g.Symbols.Names)
	return n.debug(names, "")
}

func (n *ParseTreeNode) debug(names map[int]string, indent string) string {
	result := indent + names[n.Symbol] + "\n"
	for _, child := range n.Children {
		result += child.debug(names, indent+util.Tab)
	}

	return result
}

func NewParseTreeNode(symbol int) *ParseTreeNode {
	return &ParseTreeNode{
		Symbol: symbol,
	}
}
