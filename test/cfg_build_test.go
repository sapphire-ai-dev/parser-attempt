package test

import (
	"fmt"
	"github.com/sapphire-ai-dev/parser-attempt/cfg"
	"testing"
)

func TestCFGBuild01(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_1.txt")
	fmt.Println(g.Debug())
}

func TestCFGBuild02(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_2.txt")
	fmt.Println(g.Debug())
}

func TestCFGBuild03(t *testing.T) {
	g := cfg.BuildFrom("../data/cfg_3.txt")
	fmt.Println(g.Debug())
}
