package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Config содержит конфигурацию для фильтрации
type Config struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNumber bool
	pattern    string
	filename   string
}

// MatchResult содержит информацию о совпадении
type MatchResult struct {
	line     string
	lineNum  int
	matched  bool
	filename string
}

func main() {
	config := parseFlags()
	if err := runGrep(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	config := &Config{}

	flag.IntVar(&config.after, "A", 0, "print N lines after match")
	flag.IntVar(&config.before, "B", 0, "print N lines before match")
	flag.IntVar(&config.context, "C", 0, "print N lines of context")
	flag.BoolVar(&config.count, "c", false, "print only count of matching lines")
	flag.BoolVar(&config.ignoreCase, "i", false, "ignore case")
	flag.BoolVar(&config.invert, "v", false, "invert match")
	flag.BoolVar(&config.fixed, "F", false, "treat pattern as fixed string")
	flag.BoolVar(&config.lineNumber, "n", false, "print line numbers")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] PATTERN [FILE]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	config.pattern = args[0]
	if len(args) > 1 {
		config.filename = args[1]
	}

	// Если указан контекст, устанавливаем before и after
	if config.context > 0 {
		config.before = config.context
		config.after = config.context
	}

	return config
}

func runGrep(config *Config) error {
	var scanner *bufio.Scanner
	var file *os.File

	if config.filename == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		var err error
		file, err = os.Open(config.filename)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	// Создаем регулярное выражение или фиксированную строку для поиска
	var matcher func(string) bool
	if config.fixed {
		pattern := config.pattern
		if config.ignoreCase {
			pattern = strings.ToLower(pattern)
			matcher = func(s string) bool {
				if config.ignoreCase {
					s = strings.ToLower(s)
				}
				return strings.Contains(s, pattern)
			}
		} else {
			matcher = func(s string) bool {
				return strings.Contains(s, pattern)
			}
		}
	} else {
		pattern := config.pattern
		if config.ignoreCase {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
		matcher = re.MatchString
	}

	lines, matches := readAndMatchLines(scanner, matcher, config)

	if config.count {
		fmt.Println(len(matches))
		return nil
	}

	return printResults(lines, matches, config)
}

func readAndMatchLines(scanner *bufio.Scanner, matcher func(string) bool, config *Config) ([]string, []*MatchResult) {
	var lines []string
	var matches []*MatchResult
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		matched := matcher(line)
		if config.invert {
			matched = !matched
		}

		matches = append(matches, &MatchResult{
			line:     line,
			lineNum:  lineNum,
			matched:  matched,
			filename: config.filename,
		})

		lineNum++
	}

	return lines, matches
}

func printResults(lines []string, matches []*MatchResult, config *Config) error {
	printedLines := make(map[int]bool)
	output := make([]string, 0)

	for i, match := range matches {
		if match.matched {
			// Добавляем строки до совпадения
			if config.before > 0 {
				start := max(0, i-config.before)
				for j := start; j < i; j++ {
					if !printedLines[j] {
						output = append(output, formatLine(lines[j], j+1, config, false))
						printedLines[j] = true
					}
				}
			}

			// Добавляем саму совпавшую строку
			if !printedLines[i] {
				output = append(output, formatLine(match.line, match.lineNum, config, true))
				printedLines[i] = true
			}

			// Добавляем строки после совпадения
			if config.after > 0 {
				end := min(len(matches), i+config.after+1)
				for j := i + 1; j < end; j++ {
					if !printedLines[j] {
						output = append(output, formatLine(lines[j], j+1, config, false))
						printedLines[j] = true
					}
				}
			}
		}
	}

	// Выводим все собранные строки
	for _, line := range output {
		fmt.Println(line)
	}

	return nil
}

func formatLine(line string, lineNum int, config *Config, isMatch bool) string {
	var result string

	if config.filename != "" {
		result += config.filename + ":"
	}

	if config.lineNumber {
		result += fmt.Sprintf("%d:", lineNum)
	}

	result += line

	if isMatch && config.lineNumber {
		// Добавляем маркер для совпавшей строки (как в оригинальном grep)
		result = "\033[32m" + result + "\033[0m" // Зеленый цвет для совпадения
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
