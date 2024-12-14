package main

import (
	"flag"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	buttonARegex = regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	buttonBRegex = regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	prizeRegex   = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)
)

func parseArgs() Args {
	var args Args
	flag.StringVar(&args.Input, "input", "input.txt", "")
	flag.StringVar(&args.PartFilter, "part", "", "")
	flag.BoolVar(&args.Verbose, "v", false, "verbose debug")
	flag.Parse()

	// args.Input = "input-example-1.txt"
	// args.PartFilter = "1"

	return args
}

func Part1(input string, debug *Debugger) (any, error) {
	result := 0

	var aX, aY, bX, bY, pX, pY int
	for _, line := range strings.Split(input, "\n") {
		if matches := buttonARegex.FindAllStringSubmatch(line, -1); matches != nil {
			aX, _ = strconv.Atoi(matches[0][1])
			aY, _ = strconv.Atoi(matches[0][2])
			debug.WriteFormat("A: +%d, +%d\n", aX, aY)
			continue
		}

		if matches := buttonBRegex.FindAllStringSubmatch(line, -1); matches != nil {
			bX, _ = strconv.Atoi(matches[0][1])
			bY, _ = strconv.Atoi(matches[0][2])
			debug.WriteFormat("B: +%d, +%d\n", bX, bY)
			continue
		}

		if matches := prizeRegex.FindAllStringSubmatch(line, -1); matches != nil {
			pX, _ = strconv.Atoi(matches[0][1])
			pY, _ = strconv.Atoi(matches[0][2])
			debug.WriteFormat("Prize: %d, %d\n", pX, pY)

			tokens := 0
			for b := 1; b <= 100; b++ {
				aXCalc := pX - b*bX
				if aXCalc%aX != 0 {
					continue
				}

				a := aXCalc / aX
				if bY*b+aY*a != pY {
					continue
				}

				t := a*3 + b
				if tokens == 0 || t < tokens {
					tokens = t
					debug.WriteFormat("a(%d) * 3 + b(%d) = %d\n", a, b, t)
					break
				}

				if t > tokens {
					break
				}
			}

			result += tokens
			debug.WriteFormat("Tokens: +%d => %d\n\n", tokens, result)
			continue
		}

		aX, aY, bX, bY, pX, pY = 0, 0, 0, 0, 0, 0
	}

	return result, nil
}

func Part2(input string, debug *Debugger) (any, error) {
	result := 0
	modifier := 10000000000000.0

	var aX, aY, bX, bY, pX, pY float64
	for _, line := range strings.Split(input, "\n") {
		if matches := buttonARegex.FindAllStringSubmatch(line, -1); matches != nil {
			aX, _ = strconv.ParseFloat(matches[0][1], 64)
			aY, _ = strconv.ParseFloat(matches[0][2], 64)
			debug.WriteFormat("A: +%.0f, +%.0f\n", aX, aY)
			continue
		}

		if matches := buttonBRegex.FindAllStringSubmatch(line, -1); matches != nil {
			bX, _ = strconv.ParseFloat(matches[0][1], 64)
			bY, _ = strconv.ParseFloat(matches[0][2], 64)
			debug.WriteFormat("B: +%.0f, +%.0f\n", bX, bY)
			continue
		}

		if matches := prizeRegex.FindAllStringSubmatch(line, -1); matches != nil {
			pX, _ = strconv.ParseFloat(matches[0][1], 64)
			pY, _ = strconv.ParseFloat(matches[0][2], 64)
			pX += modifier
			pY += modifier

			debug.WriteFormat("Prize: %.0f, %.0f\n", pX, pY)

			a := math.Round((pY/bY - pX/bX) / (aY/bY - aX/bX))
			b := math.Round((pX - a*aX) / bX)
			tokens := 0
			if a*aX+b*bX == pX && a*aY+b*bY == pY {
				tokens = int(a)*3 + int(b)
				result += tokens
			}
			debug.WriteFormat("Tokens: +%d => %d\n\n", tokens, result)
			continue
		}

		aX, aY, bX, bY, pX, pY = 0, 0, 0, 0, 0, 0
	}

	return result, nil
}
