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
	maxPlotWidth = 80
	padding      = 2
	usage        = "Usage: plot [-t title] [labels] , values or plot [-t title] file.csv"
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
	// check if the first row contains ANY numbers, if so, it's not a header so leave it empty
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

func main() {
	// check if arguments are passed
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	// initialize variables
	var title string
	var labels []string
	var values []float64
	var valuesArgs []string

	// check if a title flag is passed
	if os.Args[1] == "-t" {
		if len(os.Args) < 4 {
			fmt.Println(usage)
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

	// check if a csv file is passed, if so:
	if strings.HasSuffix(os.Args[1], ".csv") {
		records, _, err := readCSV(os.Args[1])
		if err != nil {
			fmt.Println("Error reading csv file:", err)
			return
		}
		if len(records) < 2 {
			fmt.Println("CSV file must have at least 2 rows")
			return
		}
		// If first column, second row is a number, use column as values
		if _, err := strconv.ParseFloat(records[1][0], 64); err == nil {

			for _, row := range records {
				value, err := strconv.ParseFloat(row[0], 64)
				if err != nil {
					fmt.Println("Error parsing value:", err)
					return
				}
				values = append(values, value)
			}
		} else {
			// use first column as labels and second column as values
			for _, row := range records {
				labels = append(labels, row[0])
				value, err := strconv.ParseFloat(row[1], 64)
				if err != nil {
					fmt.Println("Error parsing value:", err)
					fmt.Println("Either first or second column of CSV file must contain numbers")
					return
				}
				values = append(values, value)
			}
		}
	} else { // if no csv file is passed, parse arguments as labels and values
		if sepIndex != -1 {
			// parse labels and values if separator is found
			labels = os.Args[1:sepIndex]
			valuesArgs = os.Args[sepIndex+1:]
			if len(labels) != len(valuesArgs) {
				fmt.Println("Number of labels and values must match. ")
				fmt.Println("Received:", len(labels), "labels,", len(valuesArgs), "values")
				return
			}
		} else {
			// parse only values if no separator is found
			valuesArgs = os.Args[1:]
		}
	}

	// convert values to floats
	for _, arg := range valuesArgs {
		value, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Printf("invalid number: %s\n", arg)
			fmt.Println(usage)
			return
		}
		values = append(values, value)
	}

	// print the title if provided
	if title != "" {
		fmt.Println(title)
	}

	// get plot width
	termWidth, err := getTerminalWidth()
	plotWidth := min(termWidth, maxPlotWidth)
	if err != nil {
		fmt.Printf("could not determine terminal width: %s\n", err)
		return
	}

	// plot the barchart
	plot(values, labels, plotWidth)
}

// plot function to display the barchart
func plot(values []float64, labels []string, width int) {
	hasLabels := len(labels) > 0

	maxLabelLen := 0
	for _, label := range labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	// find the maximum value
	maxValue := 0.0
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}

	// find the width available for the bars
	chartWidth :=
		width -
			maxLabelLen -
			// the length of the value label, assuming rounded to int
			int(math.Log10(float64(maxValue+0.5))+1) -
			// axis and padding
			padding

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

		fmt.Print(blue)
		for j := 0; j < scaledValue; j++ {
			fmt.Print("■")
		}
		fmt.Print(reset)
		// print the value, rounded to the nearest int
		fmt.Printf(" %d\n", int(value+0.5))
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
