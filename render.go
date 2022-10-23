package main

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"go/types"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	rendCmdResp = func(cmd string, val any) {
		if val == nil {
			fmt.Println("nil")
		}

		switch {
		case matchCommand(cmd, "get"):
			renderBulkString(val)
		case matchCommand(cmd, "set"):
			renderSimpleString(val)
		case matchCommand(cmd, "info"):
			renderBulkStringDecode(val)
		case matchCommand(cmd, "hgetall"):
			renderHashPairs(val)
		case matchCommand(cmd, "memory help"):
			renderHelp(val)
		case matchCommand(cmd, "zrange"):
			renderMembers(val)
		case matchCommand(cmd, "time"):
			renderTime(val)
		case matchCommand(cmd, "lastsave"):
			renderUnixtime(val)
		case matchCommand(cmd, "lolwut"):
			renderBytes(val)
		case matchCommand(cmd, "config get"), matchCommand(cmd, "memory stats"):
			renderNestedPair(val)
		default:
			fmt.Printf("\"%v\"\n", val)
		}
	}

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
		ary := val.([]any)

		renderPair(ary, 0)
	}

	renderBytes = func(val any) {
		v := val.(string)
		fmt.Printf("%s\n", strings.TrimSpace(v))
	}

	renderUnixtime = func(val any) {
		uts := val.(int64)
		tm := time.Unix(uts, 0)

		fmt.Printf("(integer) %d\n", uts)
		fmt.Printf("(local time) %s\n", tm.Format("2006-01-02 15:04:05"))
	}

	renderSlowlog = func(val any) {

	}

	renderTime = func(val any) {
		ary := val.([]any)
		unixts, ms := ary[0].(string), ary[1].(string)
		unixTimeStamp, err := strconv.ParseInt(unixts, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		tm := time.Unix(unixTimeStamp, 0)

		fmt.Printf("(unix timestamp) %s\n", unixts)
		fmt.Printf("(millisecond) %s\n", ms)
		fmt.Printf("(convert to local timezone) %s.%s\n", tm.Format("2006-01-02 15:04:05"), ms)
	}

	renderMembers = func(val any) {
		resp := anyToString(val)
		resp = strings.Trim(resp, "[]")

		for i, v := range strings.Split(resp, " ") {
			fmt.Printf("%d) \"%v\"\n", i+1, replaceQuotes(v))
		}
	}
)

func renderPair(pairs []any, indent int) {
	var (
		keys = make([]any, 0)
		vals = make([]any, 0)
	)

	for i := 0; i < len(pairs); i += 2 {
		keys = append(keys, pairs[i])
		vals = append(vals, pairs[i+1])
	}

	for i, key := range keys {
		fmt.Printf(string(bytes.Repeat([]byte("\t"), indent)))
		fmt.Printf("%v: ", key)

		if v, ok := vals[i].([]any); ok {
			fmt.Printf("\n")
			renderPair(v, indent+1)
		} else {
			fmt.Printf("%v\n", vals[i])
		}
	}
}
