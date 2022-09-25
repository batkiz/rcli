package main

import (
	"context"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"strings"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func init() {

}

func initRedis(conf Config) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password, // no password set
		DB:       0,             // use default DB
	})
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func main() {
	rootCmd.Execute()
}

func repl() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:      conf.Prompt(),
		HistoryFile: "/tmp/readline.tmp",
		//AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	log.SetOutput(l.Stderr())

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "":
			continue
		case line == "d":
			res := confirm(l)
			fmt.Println(res)
		case line == "exit":
			goto exit
		default:
			execCmd(line)
		}
	}
exit:
}

func toAnySlice[T any](ary []T) []any {
	anyAry := make([]any, 0, len(ary))
	for _, v := range ary {
		anyAry = append(anyAry, v)
	}
	return anyAry
}

func execCmd(cmd string) {
	ary := toAnySlice(strings.Split(cmd, " "))
	redisDo(ctx, rdb, ary)
}

func redisDo(ctx context.Context, cli *redis.Client, cmd []any) {
	val, err := cli.Do(ctx, cmd...).Result()
	if err != nil {
		log.Println(err)
	} else {
		switch val.(type) {
		case []any:
			ary := val.([]any)
			for i, v := range ary {
				fmt.Printf("%d) \"%v\"\n", i+1, v)
			}
		case int:
			fmt.Printf("%d", val)
		default:
			fmt.Printf("\"%v\"\n", val)
		}
	}
}

func confirm(l *readline.Instance) bool {
	defer l.SetPrompt(l.Config.Prompt)
	l.SetPrompt("danger![y/n]")
	line, err := l.Readline()
	if err != nil {
		log.Println(err)
		return false
	}
	line = strings.ToLower(line)
	if line == "y" || line == "yes" {
		return true
	}
	return false
}
