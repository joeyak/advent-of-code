package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func main() {
	defaultInput := "input.txt"
	defaultPart := ""

	// defaultInput = "input-example-1.txt"
	// defaultPart = "2"

	var verboseDebug, outputFile bool
	var inputPath, partFilter string
	flag.StringVar(&inputPath, "input", defaultInput, "")
	flag.StringVar(&partFilter, "part", defaultPart, "")
	flag.BoolVar(&verboseDebug, "v", false, "verbose debug")
	flag.BoolVar(&outputFile, "o", false, "write output file")
	flag.Parse()

	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		slog.Error("could not read file", "path", inputPath, "err", err)
		os.Exit(1)
	}

	input := strings.TrimSuffix(strings.ReplaceAll(string(inputData), "\r\n", "\n"), "\n")

	for _, f := range []func(string, *Debugger) (any, error){Part1, Part2} {
		funcName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
		if !strings.HasSuffix(funcName, partFilter) {
			continue
		}

		debug := NewDebugBuilder(verboseDebug, fmt.Sprintf("debug-%s.txt", funcName), -1)
		defer debug.Close()

		start := time.Now()
		result, err := f(input, debug)
		duration := time.Since(start)

		if err != nil {
			slog.Error("could not run part", "func", funcName, "err", err)
			break
		}

		if outputFile {
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
		}

		debug.Close()
		slog.Info("finished running part", "func", funcName, "duration", duration, "result", result)
	}
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	// Set ID to '0'
	var ID rune = 48
	var disk []rune
	for i, r := range input {
		n := int(r) - 48

		idToAdd := ID
		if i%2 != 0 {
			idToAdd = '.'
			ID++
		}

		for i := 0; i < n; i++ {
			disk = append(disk, idToAdd)
		}
	}

	slog.Info("disk created", "length", len(disk), "lastID", ID-48)

	if ID <= '9' {
		debug.WriteString(VisualizeStep + string(disk) + "\n")
		debug.Flush()
	}

	freeIdx := 0
	numIdx := len(disk) - 1
	for {
		for disk[freeIdx] != '.' {
			freeIdx++
		}

		for disk[numIdx] == '.' {
			numIdx--
		}

		if freeIdx >= numIdx {
			break
		}

		fmt.Printf("\r%d -> %d", numIdx, freeIdx)

		disk[freeIdx] = disk[numIdx]
		disk[numIdx] = '.'

		if ID <= '9' {
			debug.WriteFormat("%s[%02d->%02d]\n%s%s\n", VisualizeStep, numIdx, freeIdx, VisualizeData, string(disk))
		} else {
			debug.WriteFormat("[%d->%d] %d\n", numIdx, freeIdx, disk[freeIdx]-48)
		}
	}

	fmt.Print("\r")

	for i, r := range disk {
		if r == '.' {
			break
		}
		result += i * int(r-48)
	}

	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0

	// Set ID to '0'
	var ID rune = 48
	var disk []rune
	for i, r := range input {
		n := int(r) - 48

		idToAdd := ID
		if i%2 != 0 {
			idToAdd = '.'
			ID++
		}

		for i := 0; i < n; i++ {
			disk = append(disk, idToAdd)
		}
	}

	slog.Info("disk created", "length", len(disk), "lastID", ID-48)

	if ID <= '9' {
		debug.WriteString(VisualizeStep + string(disk) + "\n")
		debug.Flush()
	}

	numIdx := len(disk)
	for {
		var recordedRune rune
		recordLength := 0
		for numIdx > 0 {
			numIdx--

			if disk[numIdx] == '.' && recordedRune == 0 {
				continue
			}

			if recordedRune == 0 {
				recordedRune = disk[numIdx]
			} else if disk[numIdx] != recordedRune {
				break
			}

			recordLength++
		}

		if numIdx == 0 {
			break
		}

		// Move number pointer back up one for later handling
		numIdx++

		freeIdx := 0
		freeLength := 0
		for i := range disk {
			if freeLength >= recordLength {
				freeIdx = i - freeLength
				break
			}

			if i >= numIdx {
				break
			}

			if disk[i] == '.' {
				freeLength++
			} else if freeLength != 0 {
				freeLength = 0
			}
		}

		if freeIdx == 0 {
			continue
		}

		fmt.Printf("\r[%d-%d] -> [%d,%d]           ", numIdx, numIdx+recordLength, freeIdx, freeIdx+freeLength)

		for i := 0; i < recordLength; i++ {
			disk[freeIdx+i] = recordedRune
			disk[numIdx+i] = '.'
		}

		if ID <= '9' {
			debug.WriteFormat("%s[%02d,%02d->%02d,%02d]\n%s%s\n", VisualizeStep, numIdx, numIdx+recordLength, freeIdx, freeIdx+freeLength, VisualizeData, string(disk))
			debug.Flush()
		} else {
			debug.WriteFormat("[%d,%d->%d,%d] %d\n", numIdx, numIdx+recordLength, freeIdx, freeIdx+freeLength, disk[freeIdx]-48)
		}
	}

	fmt.Print("\r")

	debug.WriteString(VisualizeEnd)
	for i, r := range disk {
		if r == '.' {
			continue
		}
		n := int(r - 48)
		result += i * n
		debug.WriteFormat("%d * %d = +%d => %d\n", i, n, i*n, result)
	}

	return result, nil
}
