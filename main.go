package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	reset        = "\033[0m"
	blue         = "\033[34m"
	bold         = "\033[1m"
	maxPlotWidth = 80
	padding      = 2
	usage        = `
Usage: 
1. plot [-t title] file.csv
2. plot [-t title] [labels ,]  values`
)

func readCSV(filename string) ([][]string, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	header := []string{}
	hasHeader := true
	for _, cell := range records[0] {
		if _, err := strconv.ParseFloat(cell, 64); err == nil {
			hasHeader = false
			break
		}
	}
	if hasHeader {
		header = records[0]
		records = records[1:]
	}

	return records, header, nil
}

func parseArgs(args []string) (string, []string, []float64, error) {
	var title string
	var labels []string
	var values []float64
	var valuesArgs []string

	if len(args) < 1 {
		return "", nil, nil, fmt.Errorf(usage)
	}

	if args[0] == "-t" {
		if len(args) < 3 {
			return "", nil, nil, fmt.Errorf(usage)
		}
		title = args[1]
		args = args[2:]
	}

	sepIndex := -1
	for i, arg := range args {
		if arg == "," {
			sepIndex = i
			break
		}
	}

	if strings.HasSuffix(args[0], ".csv") {
		records, _, err := readCSV(args[0])
		if err != nil {
			return "", nil, nil, fmt.Errorf("error reading csv file: %v", err)
		}
		if len(records) < 2 {
			return "", nil, nil, fmt.Errorf("csv file must have at least 2 rows")
		}
		if _, err := strconv.ParseFloat(records[1][0], 64); err == nil {
			for _, row := range records {
				value, err := strconv.ParseFloat(row[0], 64)
				if err != nil {
					return "", nil, nil, fmt.Errorf("error parsing value: %v", err)
				}
				values = append(values, value)
			}
		} else {
			for _, row := range records {
				labels = append(labels, row[0])
				value, err := strconv.ParseFloat(row[1], 64)
				if err != nil {
					return "", nil, nil, fmt.Errorf("error parsing value: %v", err)
				}
				values = append(values, value)
			}
		}
	} else {
		if sepIndex != -1 {
			labels = args[:sepIndex]
			valuesArgs = args[sepIndex+1:]
			if len(labels) != len(valuesArgs) {
				return "", nil, nil, fmt.Errorf("number of labels and values must match. received: %d labels, %d values", len(labels), len(valuesArgs))
			}
		} else {
			valuesArgs = args
		}
	}

	for _, arg := range valuesArgs {
		value, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return "", nil, nil, fmt.Errorf("invalid number: %s\n%s", arg, usage)
		}
		values = append(values, value)
	}

	return title, labels, values, nil
}

func plot(values []float64, labels []string, width int) {
	hasLabels := len(labels) > 0
	maxLabelLen := 0
	for _, label := range labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	maxValue := 0.0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}

	chartWidth := width - maxLabelLen - int(math.Log10(float64(maxValue+0.5))+1) - padding
	if chartWidth < 1 {
		fmt.Println("warning: labels are too long to display chart values")
	}

	scale := float64(maxValue) / float64(chartWidth)

	for i, value := range values {
		if hasLabels {
			fmt.Printf("% *s│", maxLabelLen, labels[i])
		} else {
			fmt.Print("│")
		}

		scaledValue := int(float64(value)/scale + 0.5)
		fmt.Print(blue)
		for j := 0; j < scaledValue; j++ {
			fmt.Print("■")
		}
		fmt.Print(reset)
		fmt.Printf(" %d\n", int(value+0.5))
	}
	fmt.Println()
}

func getTerminalWidth() (int, error) {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	width, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, err
	}
	return width, nil
}

func main() {
	title, labels, values, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	if title != "" {
		fmt.Print(bold)
		fmt.Println(title)
		fmt.Print(reset)
		fmt.Println()
	}

	termWidth, err := getTerminalWidth()
	plotWidth := min(termWidth, maxPlotWidth)
	if err != nil {
		fmt.Printf("could not determine terminal width: %s\n", err)
		return
	}

	plot(values, labels, plotWidth)
}
