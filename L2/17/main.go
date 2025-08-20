package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Парсинг аргументов командной строки
	timeout := flag.Int("timeout", 10, "Connection timeout in seconds")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout N] host port\n", os.Args[0])
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	// Установка соединения с таймаутом
	conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)

	// Каналы для синхронизации
	done := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для чтения из сокета и вывода в STDOUT
	go func() {
		reader := bufio.NewReader(conn)
		buf := make([]byte, 1024)

		for {
			n, err := reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("\nServer closed connection")
				} else {
					fmt.Fprintf(os.Stderr, "\nRead error: %v\n", err)
				}
				close(done)
				return
			}

			if n > 0 {
				os.Stdout.Write(buf[:n])
			}
		}
	}()

	// Горутина для чтения из STDIN и отправки в сокет
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		writer := bufio.NewWriter(conn)

		for scanner.Scan() {
			text := scanner.Text() + "\n"
			_, err := writer.WriteString(text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
				close(done)
				return
			}
			err = writer.Flush()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Flush error: %v\n", err)
				close(done)
				return
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Stdin error: %v\n", err)
		}

		close(done)
	}()

	// Ожидание завершения
	select {
	case <-done:
		fmt.Println("Connection closed")
	case <-sigCh:
		fmt.Println("\nInterrupt received, closing connection")
	}
}
