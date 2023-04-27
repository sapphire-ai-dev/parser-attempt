package cfg

import (
	"bufio"
	"github.com/sapphire-ai-dev/parser-attempt/util"
	"os"
	"strings"
)

func BuildFrom(file string) *CFG {
	result := newCFG()
	rawRules := readFile(file)

	result.Symbols.addSymbol(util.TrimWhiteSpace(rawRules[0]), false)
	var lhsList []string
	var rhsList [][]string
	for _, rawRule := range rawRules[1:] {
		lhs, rhs := interpretRule(rawRule)
		lhsList = append(lhsList, lhs)
		rhsList = append(rhsList, rhs)
	}

	nonTerminals := util.ListToSet(lhsList)
	terminals := util.KeyComplement(util.MatToSet(rhsList), nonTerminals)
	for symbol := range nonTerminals {
		result.Symbols.addSymbol(symbol, false)
	}
	for symbol := range terminals {
		result.Symbols.addSymbol(symbol, true)
	}
	if !result.Symbols.valid() {
		panic("invalid symbols")
	}

	result.Rules = make([][][]int, len(result.Symbols.Names))
	for i, lhs := range lhsList {
		result.addRule(result.Symbols.Names[lhs], util.MapOverList(result.Symbols.Names, rhsList[i]))
	}

	result.buildFirstSet()
	result.buildFollowSet()
	return result
}

func readFile(name string) []string {
	var result []string
	file, err := os.Open(name)
	util.PanicErr(err)
	defer func() { util.PrintErr(file.Close()) }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	util.PanicErr(scanner.Err())
	return result
}

func interpretRule(rule string) (string, []string) {
	sides := strings.Split(rule, "=>")
	if len(sides) != 2 {
		panic("invalid rule length")
	}

	return util.TrimWhiteSpace(sides[0]), util.TrimWhiteSpaceList(sides[1])
}
