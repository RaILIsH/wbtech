package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
)

func main() {
	// Обработка Ctrl+C
	signal.Ignore(os.Interrupt)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Minishell started. Ctrl+D to exit.")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("\nExit")
			break
		}
		if err != nil {
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		executeCommand(input)
	}
}

func executeCommand(input string) {
	// Обработка конвейеров
	if strings.Contains(input, "|") {
		executePipeline(input)
		return
	}

	// Обработка редиректов
	if strings.Contains(input, ">") || strings.Contains(input, "<") {
		executeWithRedirect(input)
		return
	}

	// Подстановка переменных окружения
	input = expandEnvironmentVariables(input)

	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	// Встроенные команды
	switch args[0] {
	case "cd":
		if len(args) > 1 {
			os.Chdir(args[1])
		}
	case "pwd":
		dir, _ := os.Getwd()
		fmt.Println(dir)
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
	case "kill":
		if len(args) > 1 {
			pid, _ := strconv.Atoi(args[1])
			process, _ := os.FindProcess(pid)
			process.Signal(os.Interrupt)
		}
	case "ps":
		cmd := exec.Command("ps")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	default:
		// Внешние команды
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

func expandEnvironmentVariables(input string) string {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			placeholder := "$" + pair[0]
			input = strings.ReplaceAll(input, placeholder, pair[1])
		}
	}
	return input
}

func executePipeline(input string) {
	commands := strings.Split(input, "|")
	var cmds []*exec.Cmd

	for _, cmdStr := range commands {
		cmdStr = strings.TrimSpace(cmdStr)
		args := strings.Fields(cmdStr)
		if len(args) == 0 {
			continue
		}
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	// Соединяем команды
	for i := 0; i < len(cmds)-1; i++ {
		stdout, _ := cmds[i].StdoutPipe()
		cmds[i+1].Stdin = stdout
	}

	// Последняя команда выводит в консоль
	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[len(cmds)-1].Stderr = os.Stderr

	// Запускаем и ждем
	for _, cmd := range cmds {
		cmd.Start()
	}
	for _, cmd := range cmds {
		cmd.Wait()
	}
}

func executeWithRedirect(input string) {
	parts := strings.Fields(input)
	var cmdArgs []string
	var inputFile, outputFile string

	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case ">":
			if i+1 < len(parts) {
				outputFile = parts[i+1]
				i++
			}
		case "<":
			if i+1 < len(parts) {
				inputFile = parts[i+1]
				i++
			}
		default:
			cmdArgs = append(cmdArgs, parts[i])
		}
	}

	if len(cmdArgs) == 0 {
		return
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	// Обработка ввода из файла
	if inputFile != "" {
		file, _ := os.Open(inputFile)
		defer file.Close()
		cmd.Stdin = file
	} else {
		cmd.Stdin = os.Stdin
	}

	// Обработка вывода в файл
	if outputFile != "" {
		file, _ := os.Create(outputFile)
		defer file.Close()
		cmd.Stdout = file
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr
	cmd.Run()
}
