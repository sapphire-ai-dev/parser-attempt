package util

import (
	"fmt"
	"strings"
)

func KeyIntersection[K, V comparable](m1, m2 map[K]V) map[K]bool {
	result := map[K]bool{}
	if m1 == nil || m2 == nil {
		return result
	}

	for k := range m1 {
		if _, seen := m2[k]; seen {
			result[k] = true
		}
	}

	return result
}

func KeyUnion[K, V comparable](m1, m2 map[K]V) map[K]bool {
	result := map[K]bool{}
	if m1 != nil {
		for k := range m1 {
			result[k] = true
		}
	}

	if m2 != nil {
		for k := range m2 {
			result[k] = true
		}
	}

	return result
}

func KeyComplement[K, V comparable](m1, m2 map[K]V) map[K]bool {
	result := map[K]bool{} // in m1 not in m2
	if m1 != nil {
		for k := range m1 {
			result[k] = true
		}
	}

	if m2 != nil {
		for k := range m2 {
			delete(result, k)
		}
	}

	return result
}

func KeyEqual[K, V comparable](m1, m2 map[K]V) bool {
	if m1 == nil && m2 == nil {
		return true
	}

	if m1 == nil || m2 == nil || len(m1) != len(m2) {
		return false
	}

	mIntersection := KeyIntersection(m1, m2)
	return len(mIntersection) == len(m1)
}

func MapCopy[K, V comparable](m map[K]V) map[K]V {
	result := map[K]V{}
	for k, v := range m {
		result[k] = v
	}

	return result
}

func MapValSet[K, V comparable](m map[K]V) map[V]bool {
	result := map[V]bool{}
	if m == nil {
		return result
	}
	for _, v := range m {
		result[v] = true
	}

	return result
}

func ListToSet[V comparable](l []V) map[V]bool {
	result := map[V]bool{}
	for _, elem := range l {
		result[elem] = true
	}

	return result
}

func MatToSet[V comparable](m [][]V) map[V]bool {
	result := map[V]bool{}
	for _, l := range m {
		for _, elem := range l {
			result[elem] = true
		}
	}

	return result
}

func MapOverList[K, V comparable](m map[K]V, l []K) []V {
	var result []V
	for _, elem := range l {
		result = append(result, m[elem])
	}

	return result
}

func InvertMap[K, V comparable](m map[K]V) map[V]K {
	result := map[V]K{}
	for k, v := range m {
		result[v] = k
	}

	return result
}

func KeyList[K, V comparable](m map[K]V) []K {
	var result []K
	for k := range m {
		result = append(result, k)
	}

	return result
}

func ContainsKey[K, V comparable](m map[K]V, k K) bool {
	_, seen := m[k]
	return seen
}

// AddedNew does not overwrite
func AddedNew[K, V comparable](m map[K]V, k K, v V) bool {
	if ContainsKey[K, V](m, k) {
		return false
	}

	m[k] = v
	return true
}

func BoolTernary[V any](b bool, t, f V) V {
	if b {
		return t
	}

	return f
}

func Max(a, b int) int {
	return BoolTernary[int](a > b, a, b)
}

func Min(a, b int) int {
	return BoolTernary[int](a < b, a, b)
}

func TrimWhiteSpace(s string) string {
	front, back := 0, len(s)
	for ; front < len(s) && s[front:front+1] == " "; front++ {
	}
	for ; back >= front && s[back-1:back] == " "; back-- {
	}
	return s[front:back]
}

func TrimWhiteSpaceList(s string) []string {
	var result []string
	if len(s) == 0 {
		return result
	}

	for _, elem := range strings.Split(s, ",") {
		result = append(result, TrimWhiteSpace(elem))
	}

	return result
}

func PrintErr(err any) {
	if err != nil {
		fmt.Println(err)
	}
}

func PanicErr(err any) {
	if err != nil {
		panic(err)
	}
}

var Tab = "  "
