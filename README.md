# kase

CLI for converting, detecting, and applying string case. Unix-friendly (pipe-safe, exit codes).

## Install

**From source (requires Go):**
```bash
go install github.com/ernilambar/kase@latest
```

**Binaries:** See [Releases](https://github.com/ernilambar/kase/releases) for Linux, macOS, and Windows.

## Usage

```bash
kase <command> <input>
kase kebab "Hello World"      # hello-world
kase snake "Hello World"     # hello_world
kase camel "hello world"     # helloWorld
kase pascal "hello world"    # HelloWorld
kase title "hello world"     # Hello World
kase detect "hello_world"    # snake
kase all "hello world"      # all forms (add --json for JSON)
```

**Flags:** `--raw` preserve accents (default: normalize to ASCII). `--json` (all only) output JSON. `-h` / `--help` for full usage.

## Behavior

Empty input → no output, exit 0. Words split on space, `_`, `-`. CamelCase boundaries detected. Accents normalized to ASCII unless `--raw`.

## License

MIT — see [LICENSE](LICENSE).
