package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	VisualizeStep = "==========STEP==========\n"
	VisualizeData = "==========DATA==========\n"
	VisualizeEnd  = "==========END==========\n"
)

type Args struct {
	Input      string
	PartFilter string
	Verbose    bool
}

func main() {
	args := parseArgs()

	inputData, err := os.ReadFile(args.Input)
	if err != nil {
		slog.Error("could not read file", "path", args.Input, "err", err)
		os.Exit(1)
	}

	input := strings.TrimSuffix(strings.ReplaceAll(string(inputData), "\r\n", "\n"), "\n")

	for _, f := range []func(string, *Debugger) (any, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
		if !strings.HasSuffix(funcName, args.PartFilter) {
			continue
		}

		debug := NewDebugBuilder(args.Verbose, fmt.Sprintf("debug-%s.txt", funcName), -1)
		defer debug.Close()

		start := time.Now()
		result, err := f(input, debug)
		duration := time.Since(start)

		if err != nil {
			slog.Error("could not run part", "func", funcName, "err", err)
			break
		}

		debug.Close()
		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

type Numbered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Abs[T Numbered](n T) T {
	if n < 0 {
		return n - n - n
	}
	return n
}

const (
	AnsiColorReset = "\033[m"

	AnsiColorDarkBlack   = "\033[0;30m"
	AnsiColorDarkRed     = "\033[0;31m"
	AnsiColorDarkGreen   = "\033[0;32m"
	AnsiColorDarkYellow  = "\033[0;33m"
	AnsiColorDarkBlue    = "\033[0;34m"
	AnsiColorDarkMagenta = "\033[0;35m"
	AnsiColorDarkCyan    = "\033[0;36m"
	AnsiColorDarkWhite   = "\033[0;37m"

	AnsiColorBlack   = "\033[0;90m"
	AnsiColorRed     = "\033[0;91m"
	AnsiColorGreen   = "\033[0;92m"
	AnsiColorYellow  = "\033[0;93m"
	AnsiColorBlue    = "\033[0;94m"
	AnsiColorMagenta = "\033[0;95m"
	AnsiColorCyan    = "\033[0;96m"
	AnsiColorWhite   = "\033[0;97m"
)

var AnsiCodes = []string{
	AnsiColorReset,
	AnsiColorDarkBlack,
	AnsiColorDarkRed,
	AnsiColorDarkGreen,
	AnsiColorDarkYellow,
	AnsiColorDarkBlue,
	AnsiColorDarkMagenta,
	AnsiColorDarkCyan,
	AnsiColorDarkWhite,
	AnsiColorBlack,
	AnsiColorRed,
	AnsiColorGreen,
	AnsiColorYellow,
	AnsiColorBlue,
	AnsiColorMagenta,
	AnsiColorCyan,
	AnsiColorWhite,
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		slog.Error("could not clear screen", "err", err)
	}
}

type Debugger struct {
	builder      bytes.Buffer
	active       bool
	filePath     string
	writeAtLen   int
	writtenBytes int
}

func NewDebugBuilder(active bool, filePath string, writeAtMB int) *Debugger {
	if writeAtMB < 0 {
		writeAtMB = 256
	}
	debugger := &Debugger{
		active:     active,
		filePath:   filePath,
		writeAtLen: writeAtMB * 1024 * 1024,
	}

	if active {
		err := debugger.write(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
		if err != nil {
			panic(fmt.Errorf("could not write file on open: %w", err))
		}
	}

	return debugger
}

func (d *Debugger) WriteString(s string) {
	if d.active {
		_, err := d.builder.WriteString(s)
		if err != nil {
			panic(fmt.Errorf("could not write string to debug builder: %w", err))
		}
		d.writeIfOverTooLarge()
	}
}

func (d *Debugger) WriteFormat(format string, a ...any) {
	if d.active {
		_, err := d.builder.WriteString(fmt.Sprintf(format, a...))
		if err != nil {
			panic(fmt.Errorf("could not write formatted string to debug builder: %w", err))
		}
		d.writeIfOverTooLarge()
	}
}

func (d *Debugger) WriteFunc(f func() string) {
	if d.active {
		_, err := d.builder.WriteString(f())
		if err != nil {
			panic(fmt.Errorf("could not write formatted string to debug builder: %w", err))
		}
		d.writeIfOverTooLarge()
	}
}

func (d *Debugger) writeIfOverTooLarge() {
	if d.builder.Len() > d.writeAtLen {
		err := d.write(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
		if err != nil {
			panic(fmt.Errorf("could not write periodic data to file: %w", err))
		}
	}
}

func (d *Debugger) Close() {
	if d.builder.Len() == 0 {
		return
	}

	err := d.write(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
	if err != nil {
		panic(fmt.Errorf("could not write file on close: %w", err))
	}
}

func (d *Debugger) Len() int {
	return d.builder.Len() + d.writtenBytes
}

func (d *Debugger) Flush() {
	if d.builder.Len() > 0 {
		err := d.write(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
		if err != nil {
			panic(fmt.Errorf("could not write periodic data to file: %w", err))
		}
	}
}

func (d *Debugger) write(osFlags int) error {
	file, err := os.OpenFile(d.filePath, osFlags, 0644)
	if err != nil {
		return fmt.Errorf("could not open file for write: %w", err)
	}
	defer file.Close()

	data := d.builder.Bytes()
	for _, code := range AnsiCodes {
		data = bytes.ReplaceAll(data, []byte(code), nil)
	}

	n, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}
	d.writtenBytes += n

	d.builder.Reset()
	return nil
}

type Hash[T comparable] map[T]struct{}

func (h Hash[T]) Add(k T) {
	h[k] = struct{}{}
}
