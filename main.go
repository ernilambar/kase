package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ernilambar/kase/internal/caseconv"
	"github.com/urfave/cli/v2"
)

// Version is set at build time via -ldflags "-X main.Version=..."
var Version = "0.1.0"

func main() {
	app := &cli.App{
		Name:    "kase",
		Usage:   "convert string case",
		Version: Version,
		Description: `Convert strings between kebab, snake, camel, pascal, title case.
If input is omitted, read from stdin (e.g. cat file.txt | kase kebab).`,
		Commands: []*cli.Command{
			{
				Name:   "kebab",
				Usage:  "lowercase, hyphen-separated (hello-world)",
				Flags:  sharedFlags(),
				Action: runSingle(caseconv.ToKebab),
			},
			{
				Name:   "snake",
				Usage:  "lowercase, underscore-separated (hello_world)",
				Flags:  sharedFlags(),
				Action: runSingle(caseconv.ToSnake),
			},
			{
				Name:   "camel",
				Usage:  "camelCase (helloWorld)",
				Flags:  sharedFlags(),
				Action: runSingle(caseconv.ToCamel),
			},
			{
				Name:   "pascal",
				Usage:  "PascalCase (HelloWorld)",
				Flags:  sharedFlags(),
				Action: runSingle(caseconv.ToPascal),
			},
			{
				Name:   "title",
				Usage:  "Title Case (Hello World)",
				Flags:  sharedFlags(),
				Action: runSingle(caseconv.ToTitle),
			},
			{
				Name:   "detect",
				Usage:  "detect case of input string",
				Flags:  sharedFlags(),
				Action: runSingle(detectOnly),
			},
			{
				Name:  "all",
				Usage: "output all case conversions",
				Flags: append(sharedFlags(),
					&cli.BoolFlag{
						Name:  "json",
						Usage: "output all conversions as JSON",
					},
				),
				Action: runAllCmd,
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Fprint(os.Stderr, "kase: command required\n\n")
			_ = cli.ShowAppHelp(cCtx)
			os.Exit(1)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func sharedFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "raw",
			Usage: "preserve accented characters (default: normalize to ASCII)",
		},
	}
}

func getInput(cCtx *cli.Context) string {
	input := cCtx.Args().First()
	if input == "" {
		input = readStdin()
	}
	return input
}

func runSingle(fn func(string, bool) string) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		input := getInput(cCtx)
		if strings.TrimSpace(input) == "" {
			return nil
		}
		raw := cCtx.Bool("raw")
		result := fn(input, raw)
		if result != "" {
			fmt.Println(result)
		}
		return nil
	}
}

func detectOnly(s string, _ bool) string { return caseconv.Detect(s) }

func runAllCmd(cCtx *cli.Context) error {
	input := getInput(cCtx)
	if strings.TrimSpace(input) == "" {
		return nil
	}
	raw := cCtx.Bool("raw")
	jsonOutput := cCtx.Bool("json")

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
		return enc.Encode(out)
	}
	fmt.Printf("kebab  : %s\n", kebab)
	fmt.Printf("snake  : %s\n", snake)
	fmt.Printf("camel  : %s\n", camel)
	fmt.Printf("pascal : %s\n", pascal)
	fmt.Printf("title  : %s\n", title)
	return nil
}

// readStdin reads all of stdin and returns it as one string. Newlines are normalized to spaces.
func readStdin() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return ""
	}
	s := string(data)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}
