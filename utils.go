package main

import (
	"fmt"
	"strings"
)

func toAnySlice[T any](ary []T) []any {
	anyAry := make([]any, 0, len(ary))
	for _, v := range ary {
		anyAry = append(anyAry, v)
	}
	return anyAry
}

func removeQuotes(s string) string {
	x := strings.TrimLeft(s, `"`)
	x = strings.TrimRight(s, `"`)

	return x
}

func anyToString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	default:
		return fmt.Sprint(val)
	}
}

func doubleQuotes(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func matchCommand(s, cmd string) bool {
	return strings.HasPrefix(strings.ToLower(s), cmd)
}

func stringToArgsSlice(cmd string) []any {
	ary := strings.Split(cmd, " ")
	args := make([]any, 0, len(ary))
	for _, s := range ary {
		args = append(args, strings.Trim(s, `"'`))
	}

	return args
}

func replaceQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
