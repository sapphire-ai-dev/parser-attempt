package test

import (
	"github.com/sapphire-ai-dev/parser-attempt/cfg"
	"github.com/sapphire-ai-dev/parser-attempt/lr"
	"testing"
)

func TestSLRBuild01(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_1.txt")
	lr.BuildSLRParser(g)
}

func TestSLRBuild02(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_2.txt")
	lr.BuildSLRParser(g)
}

func TestSLRBuild03(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_3.txt")
	lr.BuildSLRParser(g)
}
