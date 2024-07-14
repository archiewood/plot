package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	// check if arguments are passed
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go [-t title] [labels] , [values] or go run main.go [-t title] [values]")
		return
	}

	// initialize variables
	var title string
	var labels []string
	var values []int
	var valuesArgs []string

	// check if a title flag is passed
	if os.Args[1] == "-t" {
		if len(os.Args) < 4 {
			fmt.Println("usage: go run main.go -t [title] [labels] , [values] or go run main.go -t [title] [values]")
			return
		}
		title = os.Args[2]
		os.Args = append(os.Args[:1], os.Args[3:]...)
	}

	// find the separator index
	sepIndex := -1
	for i, arg := range os.Args {
		if arg == "," {
			sepIndex = i
			break
		}
	}

	if sepIndex != -1 {
		// parse labels and values if separator is found
		labels = os.Args[1:sepIndex]
		valuesArgs = os.Args[sepIndex+1:]
		if len(labels) != len(valuesArgs) {
			fmt.Println("the number of labels and values must be the same")
			return
		}
	} else {
		// parse only values if no separator is found
		valuesArgs = os.Args[1:]
	}

	// convert values to integers
	for _, arg := range valuesArgs {
		value, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("invalid number: %s\n", arg)
			fmt.Println("usage: go run main.go [labels] , [values] or go run main.go [values]")
			return
		}
		values = append(values, value)
	}

	// print the title if provided
	if title != "" {
		fmt.Println(title)
	}

	// get terminal width
	termWidth, err := getTerminalWidth()
	if err != nil {
		fmt.Printf("could not determine terminal width: %s\n", err)
		return
	}

	// plot the barchart
	plot(values, labels, termWidth)
}

// plot function to display the barchart
func plot(values []int, labels []string, width int) {
	hasLabels := len(labels) > 0

	maxLabelLen := 0
	for _, label := range labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	// find the maximum value
	maxValue := 0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}

	// find the width available for the bars
	chartWidth :=
		width -
			maxLabelLen -
			// the length of the value label
			int(math.Log10(float64(maxValue))+0.5) -
			// the padding
			2

	fmt.Print("termWidth:", width, "=\n")
	for j := 0; j < width; j++ {
		fmt.Print("■")
	}

	// calculate the scale factor for the values by dividing the maximum value by the terminal width, accounting for the label length and the value label length
	scale := float64(maxValue) / float64(chartWidth)

	for i, value := range values {
		if hasLabels {
			label := labels[i]
			// print the label, using the max label length
			fmt.Printf("% *s│", maxLabelLen, label)
		} else {
			fmt.Print("│")
		}

		// calculate the number of blocks to, rounding the value to the nearest integer
		scaledValue := int(float64(value)/scale + 0.5)

		for j := 0; j < scaledValue; j++ {
			fmt.Print("■")
		}
		// print the value
		fmt.Printf(" %d\n", value)
	}
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
