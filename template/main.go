package main

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	_ "embed"
)

//go:embed input.txt
var input string

func main() {
	for _, f := range []func() (any, string, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]

		start := time.Now()
		result, debug, err := f()
		duration := time.Since(start)

		if err != nil {
			slog.Error("could not run part", "func", funcName, "err", err)
			break
		}

		file, err := os.OpenFile(fmt.Sprintf("output-%s.txt", funcName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		if err != nil {
			slog.Error("could not create or append output file", "err", err)
			os.Exit(1)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("%v", result))
		if err != nil {
			slog.Error("could not write result", "func", funcName, "result", result, "err", err)
			break
		}

		debugFile, err := os.OpenFile(fmt.Sprintf("debug-%s.txt", funcName), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			slog.Error("could not create or append debug file", "err", err)
			os.Exit(1)
		}
		defer debugFile.Close()

		_, err = debugFile.WriteString(debug)
		if err != nil {
			slog.Error("could not write debug", "func", funcName, "err", err)
			break
		}

		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1() (any, string, error) {
	return "", "", nil
}

func Part2() (any, string, error) {
	return "", "", nil
}