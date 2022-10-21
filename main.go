package main

//go:generate go run data/generate.go > const.go

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
		Prompt:          conf.Prompt(),
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
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
		case line == "exit":
			goto exit
		default:
			if confirm(l, line) {
				execCmd(line)
			}
		}
	}
exit:
}

func execCmd(cmd string) {
	//ary := toAnySlice(strings.Split(cmd, " "))
	redisDo(ctx, rdb, cmd)
}

func redisDo(ctx context.Context, cli *redis.Client, cmd string) {
	cmdAry := stringToArgsSlice(cmd)
	val, err := cli.Do(ctx, cmdAry...).Result()

	if err != nil {
		log.Println(err)
	} else {
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
		default:
			fmt.Printf("\"%v\"\n", val)
		}
	}
}

func confirm(l *readline.Instance, line string) bool {
	defer l.SetPrompt(l.Config.Prompt)

	cmdAry := strings.Split(line, " ")
	cmd := strings.ToUpper(cmdAry[0])

	msg, ok := dangerousCommands[cmd]
	if !ok {
		return true
	}

	alert := msg + " [y/n]"

	l.SetPrompt(alert)
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
