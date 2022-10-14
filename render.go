package main

import (
	"fmt"
	"github.com/gookit/color"
	"go/types"
	"strings"
)

var (
	renderSimpleString = func(val any) {
		fmt.Printf("%v\n", val)
	}

	renderInt = func(val any) {
		color.Note.Print("(integer) ")
		fmt.Printf("%d\n", val)
	}

	renderList = func(val any) {
		ary := val.([]any)
		for i, v := range ary {
			fmt.Printf("%d) \"%v\"\n", i+1, replaceQuotes(anyToString(v)))
		}
	}

	renderBulkStringDecode = func(val any) {
		s := removeQuotes(anyToString(val))
		fmt.Println(strings.TrimRight(s, "\n"))
	}

	renderBulkString = func(val any) {
		fmt.Printf("\"%v\"\n", val)
	}

	renderStringOrInt = func(val any) {
		switch val.(type) {
		case int:
			renderInt(val)
		default:
			renderBulkString(val)
		}
	}

	renderListOrString = func(val any) {
		switch val.(type) {
		case types.Array:
			renderList(val)
		default:
			renderBulkString(val)
		}
	}

	renderHashPairs = func(val any) {
		hashes := func() []string {
			x := strings.TrimLeft(anyToString(val), "[")
			x = strings.TrimRight(x, "]")
			return strings.Split(x, " ")
		}()

		hashPairs := make([][2]string, 0, len(hashes)/2)

		for i := 0; i < len(hashes); i += 2 {
			pair := [2]string{
				hashes[i],
				hashes[i+1],
			}
			hashPairs = append(hashPairs, pair)
		}

		sb := strings.Builder{}

		for i, pair := range hashPairs {
			sb.WriteString(
				fmt.Sprintf(
					"%d) %s\n   %s\n",
					i+1,
					doubleQuotes(pair[0]),
					doubleQuotes(pair[1]),
				),
			)
		}

		fmt.Printf(sb.String())
	}

	renderSubscribe = func(val any) {

	}

	renderHelp = func(val any) {
		ary := val.([]any)

		sb := strings.Builder{}
		for _, s := range ary {
			sb.WriteString(anyToString(s))
			sb.WriteString("\n")
		}

		fmt.Printf(sb.String())
	}
	renderNestedPair = func(val any) {

	}

	renderBytes = func(val any) {

	}

	renderUnixtime = func(val any) {

	}

	renderSlowlog = func(val any) {

	}

	renderTime = func(val any) {

	}

	renderMembers = func(val any) {
		resp := anyToString(val)
		resp = strings.Trim(resp, "[]")

		for i, v := range strings.Split(resp, " ") {
			fmt.Printf("%d) \"%v\"\n", i+1, replaceQuotes(v))
		}
	}
)
