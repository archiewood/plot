package main

import (
	"fmt"
	"os"
	"strconv"
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

	// plot the barchart
	plot(values, labels)
}

// plot function to display the barchart
func plot(values []int, labels []string) {
	hasLabels := len(labels) > 0

	maxLabelLen := 0
	for _, label := range labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	for i, value := range values {
		if hasLabels {
			label := labels[i]
			// print the label, using the max label length
			fmt.Printf("% *s │", maxLabelLen, label)
		} else {
			fmt.Print("│")
		}
		// print the bar
		for j := 0; j < value; j++ {
			fmt.Print("■")
		}
		// print the value
		fmt.Printf(" %d\n", value)
	}
}
