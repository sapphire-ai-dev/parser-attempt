package test

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/cfg"
	"github.com/sapphire-ai-dev/parser-attempt/lr"
	"strings"
	"testing"
)

func TestSLRParse01(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_1.txt")
	p := lr.BuildSLRParser(g)
	fmt.Println(p.Debug())
	n := p.ParseRaw(strings.Split("b a a", " "))
	fmt.Println(n.Debug(g))
}

func TestSLRParse02(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_2.txt")
	p := lr.BuildSLRParser(g)
	fmt.Println(p.Debug())
	n := p.ParseRaw(strings.Split("int + int * ( int + int )", " "))
	fmt.Println(n.Debug(g))
}

func TestSLRParse03(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_3.txt")
	p := lr.BuildSLRParser(g)
	fmt.Println(p.Debug())
	n := p.ParseRaw(strings.Split("int + int * ( int + int )", " "))
	fmt.Println(n.Debug(g))
}
