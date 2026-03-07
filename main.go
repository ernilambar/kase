package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ernilambar/kase/internal/caseconv"
)

const usage = `kase - convert string case

Usage:
  kase <command> <input>

Commands:
  kebab   lowercase, hyphen-separated (hello-world)
  snake   lowercase, underscore-separated (hello_world)
  camel   camelCase (helloWorld)
  pascal  PascalCase (HelloWorld)
  title   Title Case (Hello World)
  detect  detect case of input string
  all     output all case conversions

Flags:
  --json  (all only) output all conversions as JSON
  --raw   preserve accented characters (default: normalize to ASCII)

Example:
  kase kebab "Hello World"
  kase kebab --raw "Planeta Envíos"
  kase all "hello world"
  kase all --json "hello world"
`

var validCommands = map[string]bool{
	"kebab":  true,
	"snake":  true,
	"camel":  true,
	"pascal": true,
	"title":  true,
	"detect": true,
	"all":    true,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	arg1 := os.Args[1]
	if arg1 == "-h" || arg1 == "--help" {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := arg1

	// "all" accepts: kase all [--json] [--raw] <input>
	// Other commands: kase <cmd> [--raw] <input>
	var input string
	var allJSON bool
	var raw bool
	if cmd == "all" {
		if len(os.Args) < 3 {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(1)
		}
		for i := 2; i < len(os.Args)-1; i++ {
			switch os.Args[i] {
			case "--json":
				allJSON = true
			case "--raw":
				raw = true
			}
		}
		input = os.Args[len(os.Args)-1]
		if input == "--raw" || input == "--json" {
			fmt.Fprint(os.Stderr, "kase: all requires an input string\n\n"+usage)
			os.Exit(1)
		}
	} else {
		if len(os.Args) < 3 {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(1)
		}
		if os.Args[2] == "--raw" {
			if len(os.Args) < 4 {
				fmt.Fprint(os.Stderr, "kase: --raw requires an input string\n\n"+usage)
				os.Exit(1)
			}
			raw = true
			input = os.Args[3]
		} else {
			input = os.Args[2]
		}
	}

	if !validCommands[cmd] {
		fmt.Fprintf(os.Stderr, "kase: unknown command %q\n\n%s", cmd, usage)
		os.Exit(1)
	}

	// Empty or whitespace-only input: print nothing, exit 0 (script-friendly)
	if strings.TrimSpace(input) == "" {
		os.Exit(0)
	}

	if cmd == "all" {
		runAll(input, allJSON, raw)
		os.Exit(0)
	}

	var result string
	switch cmd {
	case "kebab":
		result = caseconv.ToKebab(input, raw)
	case "snake":
		result = caseconv.ToSnake(input, raw)
	case "camel":
		result = caseconv.ToCamel(input, raw)
	case "pascal":
		result = caseconv.ToPascal(input, raw)
	case "title":
		result = caseconv.ToTitle(input, raw)
	case "detect":
		result = caseconv.Detect(input)
	default:
		result = input
	}

	if result != "" {
		fmt.Println(result)
	}
	os.Exit(0)
}

func runAll(input string, jsonOutput bool, raw bool) {
	kebab := caseconv.ToKebab(input, raw)
	snake := caseconv.ToSnake(input, raw)
	camel := caseconv.ToCamel(input, raw)
	pascal := caseconv.ToPascal(input, raw)
	title := caseconv.ToTitle(input, raw)

	if jsonOutput {
		out := struct {
			Kebab  string `json:"kebab"`
			Snake  string `json:"snake"`
			Camel  string `json:"camel"`
			Pascal string `json:"pascal"`
			Title  string `json:"title"`
		}{kebab, snake, camel, pascal, title}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(out)
		return
	}
	fmt.Printf("kebab  : %s\n", kebab)
	fmt.Printf("snake  : %s\n", snake)
	fmt.Printf("camel  : %s\n", camel)
	fmt.Printf("pascal : %s\n", pascal)
	fmt.Printf("title  : %s\n", title)
}
