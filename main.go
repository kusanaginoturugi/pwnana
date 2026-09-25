package main

import (
	crand "crypto/rand"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

//go:embed goroawase.json
var goroawaseJSON []byte

const (
	defaultLength = 16
	minLength     = 8
)

var (
	layout = struct {
		leftConsonants  string
		leftVowels      string
		rightConsonants string
		rightVowels     string
	}{
		leftConsonants:  "wrtsfgzxcvbq",
		leftVowels:      "ea",
		rightConsonants: "yphjklnm",
		rightVowels:     "uio",
	}

	digraphs    = []string{"oo", "ee", "ii", "ou", "ai", "au"}
	leet        = map[byte]byte{'i': '!', 'a': '@', 's': '$', 'e': '3', 'o': '0', 't': '7', 'q': '9', 'b': '6', 'g': '8'}
	wordSymbols = "-+%~&/_"
	bannedUpper = map[byte]bool{'O': true, 'I': true}
	vowels      = "aeiou"
	errHelp     = errors.New("help requested")
)

type goro struct {
	label  string
	digits string
}

type options struct {
	length      int
	count       int
	countSet    bool
	digit       bool
	symbol      bool
	uppercase   bool
	wordSymbols bool
	goroawase   bool
}

func usage() string {
	return `Usage: pwnana [length] [count] [-d|--digit] [-s|--symbol] [--word-symbols] [-u|--uppercase] [-g|--goroawase] [-h|--help]

Pronounceable, keyboard-alternating password generator

Arguments:
  length           Password length (min 8, default 16)
  count            Number of passwords to generate (default: about half the terminal height)

Options:
  -h, --help       Show this help message and exit
  -d, --digit      Include one digit
  -s, --symbol     Include one symbol
  -u, --uppercase  Include one uppercase letter
      --word-symbols
                  Use terminal word-selectable symbols (-+%~&/_)
  -g, --goroawase  Embed goroawase digits with hint`
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if errors.Is(err, errHelp) {
		fmt.Println(usage())
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwnana: %v\n\n%s\n", err, usage())
		os.Exit(2)
	}

	goroawase, err := loadGoroawase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwnana: %v\n", err)
		os.Exit(1)
	}

	count := opts.count
	if !opts.countSet {
		count = defaultCount()
	}
	symbol := opts.symbol || opts.wordSymbols

	for i := 0; i < count; i++ {
		pw, hint, err := generate(opts.length, opts.digit, symbol, opts.goroawase, opts.uppercase, opts.wordSymbols, goroawase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pwnana: %v\n", err)
			os.Exit(1)
		}
		if hint != "" {
			fmt.Printf("%s  [%s]\n", pw, hint)
		} else {
			fmt.Println(pw)
		}
	}
}

func parseArgs(args []string) (options, error) {
	opts := options{length: defaultLength}
	positionals := []string{}
	onlyPositionals := false

	for _, arg := range args {
		if onlyPositionals {
			positionals = append(positionals, arg)
			continue
		}
		if arg == "--" {
			onlyPositionals = true
			continue
		}
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positionals = append(positionals, arg)
			continue
		}

		switch arg {
		case "-h", "--help":
			return opts, errHelp
		case "-d", "--digit":
			opts.digit = true
		case "-s", "--symbol":
			opts.symbol = true
		case "-u", "--uppercase":
			opts.uppercase = true
		case "-g", "--goroawase":
			opts.goroawase = true
		case "--word-symbols":
			opts.wordSymbols = true
		default:
			if strings.HasPrefix(arg, "--") {
				return opts, fmt.Errorf("unknown option %s", arg)
			}
			if err := parseShortCluster(arg[1:], &opts); err != nil {
				return opts, err
			}
		}
	}

	if len(positionals) > 2 {
		return opts, fmt.Errorf("too many positional arguments")
	}
	if len(positionals) >= 1 {
		length, err := strconv.Atoi(positionals[0])
		if err != nil {
			return opts, fmt.Errorf("invalid length %q", positionals[0])
		}
		opts.length = max(minLength, length)
	}
	if len(positionals) == 2 {
		count, err := strconv.Atoi(positionals[1])
		if err != nil {
			return opts, fmt.Errorf("invalid count %q", positionals[1])
		}
		opts.count = evenCount(count)
		opts.countSet = true
	}
	if len(positionals) == 0 {
		opts.length = defaultLength
	}
	return opts, nil
}

func parseShortCluster(cluster string, opts *options) error {
	if cluster == "" {
		return fmt.Errorf("unknown option -")
	}
	for _, c := range cluster {
		switch c {
		case 'h':
			return errHelp
		case 'd':
			opts.digit = true
		case 's':
			opts.symbol = true
		case 'u':
			opts.uppercase = true
		case 'g':
			opts.goroawase = true
		default:
			return fmt.Errorf("unknown option -%c", c)
		}
	}
	return nil
}

func loadGoroawase() ([]goro, error) {
	var raw [][]string
	if err := json.Unmarshal(goroawaseJSON, &raw); err != nil {
		return nil, fmt.Errorf("load goroawase.json: %w", err)
	}

	result := make([]goro, 0, len(raw))
	for _, item := range raw {
		if len(item) != 2 || len(item[1]) != 4 {
			continue
		}
		result = append(result, goro{label: item[0], digits: item[1]})
	}
	return result, nil
}

func generate(length int, digit, symbol, useGoro, uppercase, useWordSymbols bool, goroawase []goro) (string, string, error) {
	chunks := []string{}
	hand := "left"
	if randomInt(2) == 1 {
		hand = "right"
	}

	for totalLen(chunks) < length {
		if randomInt(10) < 2 {
			chunks = append(chunks, digraphs[randomInt(len(digraphs))])
			continue
		}
		unit, nextHand := makeUnit(hand)
		chunks = append(chunks, unit)
		hand = nextHand
	}

	result := []byte(strings.Join(chunks, "")[:length])
	hint := ""
	usedPositions := map[int]bool{}
	insertionPositions := fixedPositions(length)

	if useGoro && len(goroawase) > 0 {
		item := goroawase[randomInt(len(goroawase))]
		pos := pickInt(fixedGoroStartPositions(length))
		for i := 0; i < len(item.digits) && pos+i < length; i++ {
			result[pos+i] = item.digits[i]
			usedPositions[pos+i] = true
		}
		hint = item.label
	} else if digit && !applyLeet(result, true) {
		pos := pickInsertionPosition(insertionPositions, usedPositions)
		result[pos] = byte('0' + randomInt(10))
		usedPositions[pos] = true
	}

	if symbol && !useWordSymbols && !applyLeet(result, false) {
		symbols := []byte{}
		for _, v := range leet {
			if !isDigit(v) {
				symbols = append(symbols, v)
			}
		}
		pos := pickInsertionPosition(insertionPositions, usedPositions)
		result[pos] = symbols[randomInt(len(symbols))]
		usedPositions[pos] = true
	} else if symbol && useWordSymbols {
		pos := pickInsertionPosition(insertionPositions, usedPositions)
		result[pos] = wordSymbols[randomInt(len(wordSymbols))]
		usedPositions[pos] = true
	}

	last := result[len(result)-1]
	if isConsonant(last) && lower(last) != 'n' {
		allVowels := layout.leftVowels + layout.rightVowels
		result[len(result)-1] = allVowels[randomInt(len(allVowels))]
	}

	applyDefaultRules(result)
	if uppercase {
		applyUppercase(result)
	}
	return string(result), hint, nil
}

func makeUnit(hand string) (string, string) {
	if hand == "left" {
		if randomInt(2) == 1 {
			return string(layout.rightVowels[randomInt(len(layout.rightVowels))]) +
				string(layout.leftConsonants[randomInt(len(layout.leftConsonants))]), "left"
		}
		return string(layout.rightConsonants[randomInt(len(layout.rightConsonants))]) +
			string(layout.leftVowels[randomInt(len(layout.leftVowels))]), "left"
	}

	if randomInt(2) == 1 {
		return string(layout.leftVowels[randomInt(len(layout.leftVowels))]) +
			string(layout.rightConsonants[randomInt(len(layout.rightConsonants))]), "right"
	}
	return string(layout.leftConsonants[randomInt(len(layout.leftConsonants))]) +
		string(layout.rightVowels[randomInt(len(layout.rightVowels))]), "right"
}

func applyLeet(result []byte, wantDigit bool) bool {
	candidates := []int{}
	for i, c := range result {
		replacement, ok := leet[c]
		if !ok {
			continue
		}
		if isDigit(replacement) == wantDigit {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return false
	}
	idx := pickInt(candidates)
	result[idx] = leet[result[idx]]
	return true
}

func applyDefaultRules(result []byte) {
	for i, c := range result {
		if bannedUpper[c] {
			result[i] = lower(c)
		}
	}
}

func applyUppercase(result []byte) bool {
	candidates := []int{}
	for i, c := range result {
		if isAlpha(c) && !bannedUpper[upper(c)] {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		return false
	}
	idx := pickInt(candidates)
	result[idx] = upper(result[idx])
	return true
}

func fixedPositions(length int) []int {
	positions := []int{}
	for i := 3; i < length-1; i += 4 {
		if i >= 1 && i <= length-2 {
			positions = append(positions, i)
		}
	}
	if len(positions) == 0 {
		return []int{max(1, length/2)}
	}
	return positions
}

func fixedGoroStartPositions(length int) []int {
	positions := []int{}
	for i := 0; i < length-3; i += 4 {
		positions = append(positions, i)
	}
	if len(positions) == 0 {
		return []int{0}
	}
	return positions
}

func pickInsertionPosition(positions []int, used map[int]bool) int {
	candidates := []int{}
	for _, p := range positions {
		if !used[p] {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) > 0 {
		return pickInt(candidates)
	}
	return pickInt(positions)
}

func pickInt(values []int) int {
	return values[randomInt(len(values))]
}

func totalLen(chunks []string) int {
	total := 0
	for _, c := range chunks {
		total += len(c)
	}
	return total
}

func evenCount(count int) int {
	count = max(2, count)
	if count%2 != 0 {
		count++
	}
	return count
}

func defaultCount() int {
	return evenCount(terminalLines() / 2)
}

func terminalLinesFallback() int {
	if lines, err := strconv.Atoi(os.Getenv("LINES")); err == nil && lines > 0 {
		return lines
	}
	return 24
}

func randomInt(n int) int {
	if n <= 0 {
		panic("randomInt called with non-positive bound")
	}
	value, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(value.Int64())
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isConsonant(c byte) bool {
	return isAlpha(c) && !strings.ContainsRune(vowels, rune(lower(c)))
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - ('a' - 'A')
	}
	return c
}
