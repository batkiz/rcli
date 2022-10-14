package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/jszwec/csvutil"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

type DangerousCommand struct {
	Command string `csv:"Command"`
	Reason  string `csv:"Reason"`
}

type Command struct {
	Group    string `csv:"Group"`
	Command  string `csv:"Command"`
	Syntax   string `csv:"Syntax"`
	Callback string `csv:"Callback"`
}

var (
	//go:embed dangerous_commands.csv
	dangerousCommandsCsv []byte

	//go:embed command_syntax.csv
	commandSyntaxCsv []byte

	//go:embed const.tmpl
	tmpl string
)

// 生成危险命令的提示
func genDangerCommands() string {
	var commands []DangerousCommand
	if err := csvutil.Unmarshal(dangerousCommandsCsv, &commands); err != nil {
		log.Fatalln(err)
	}

	res := strings.Builder{}

	for _, cmd := range commands {
		res.WriteString(fmt.Sprintf("\"%s\": \"%s\",\n", cmd.Command, cmd.Reason))
	}

	return res.String()
}

// 生成命令的补全
func genCommands() string {
	var commands []Command
	if err := csvutil.Unmarshal(commandSyntaxCsv, &commands); err != nil {
		log.Fatalln(err)
	}

	var (
		SingleCmd = make([]string, 0, len(commands))
		MultiCmd  = make(map[string][]string)
		res       = strings.Builder{}
	)

	for _, cmd := range commands {
		if strings.Contains(cmd.Command, " ") {
			ary := strings.Split(cmd.Command, " ")
			if v, ok := MultiCmd[ary[0]]; ok {
				MultiCmd[ary[0]] = append(v, ary[1])
			} else {
				MultiCmd[ary[0]] = []string{ary[1]}
			}
		} else {
			SingleCmd = append(SingleCmd, cmd.Command)
		}
	}

	// 大写
	for _, s := range SingleCmd {
		res.WriteString(fmt.Sprintf("readline.PcItem(\"%s\"),\n", s))
	}
	//	小写
	for _, s := range SingleCmd {
		res.WriteString(fmt.Sprintf("readline.PcItem(\"%s\"),\n", strings.ToLower(s)))
	}

	// 大写
	for k, v := range MultiCmd {
		res.WriteString("readline.PcItem(\n")
		res.WriteString(fmt.Sprintf("\"%s\",\n", k))
		for _, s := range v {
			res.WriteString(fmt.Sprintf("readline.PcItem(\"%s\"),\n", s))
		}
		res.WriteString("),\n")
	}

	// 小写
	for k, v := range MultiCmd {
		res.WriteString("readline.PcItem(\n")
		res.WriteString(fmt.Sprintf("\"%s\",\n", strings.ToLower(k)))
		for _, s := range v {
			res.WriteString(fmt.Sprintf("readline.PcItem(\"%s\"),\n", strings.ToLower(s)))
		}
		res.WriteString("),\n")
	}

	return res.String()
}

type TmplData struct {
	Commands       string
	DangerCommands string
}

func main() {
	b, err := execTmpl(TmplData{
		Commands:       genCommands(),
		DangerCommands: genDangerCommands(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	formattedContent, err := format.Source(b)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("const.go", formattedContent, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func execTmpl(data TmplData) ([]byte, error) {
	t := template.Must(template.New("tmpl").Parse(tmpl))

	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "tmpl", data)

	if err != nil {
		log.Println("executing template:", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
