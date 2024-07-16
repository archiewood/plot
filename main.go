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
	reset         = "\033[0m"
	blue          = "\033[34m"
	bold          = "\033[1m"
	maxPlotWidth  = 800
	maxPlotHeight = 20
	paddingWidth  = 2
	paddingHeight = 2
	usage         = `
Usage: 
1. plot [-t title] [-c chartType] file.csv
2. plot [-t title] [-c chartType] [labels ,]  values

Options:
  -t title 	 Title of the chart
  -c chartType	 Type of chart (bar, column)`
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

func parseArgs(args []string) (string, string, []string, []float64, error) {
	var title string
	var chartType string
	var labels []string
	var values []float64
	var valuesArgs []string

	if len(args) < 1 {
		return "", "", nil, nil, fmt.Errorf(usage)
	}

	i := 0

	for i < len(args) {
		switch args[i] {
		case "-t":
			if i+1 >= len(args) {
				return "", "", nil, nil, fmt.Errorf(usage)
			}
			title = args[i+1]
			i += 2
		case "-c":
			if i+1 >= len(args) {
				return "", "", nil, nil, fmt.Errorf(usage)
			}
			chartType = args[i+1]
			i += 2
		default:
			args = args[i:]
			i = len(args) // end loop
		}
	}

	if len(args) == 0 {
		return "", "", nil, nil, fmt.Errorf(usage)
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
			return "", "", nil, nil, fmt.Errorf("error reading csv file: %v", err)
		}
		if len(records) < 2 {
			return "", "", nil, nil, fmt.Errorf("csv file must have at least 2 rows")
		}
		if _, err := strconv.ParseFloat(records[1][0], 64); err == nil {
			for _, row := range records {
				value, err := strconv.ParseFloat(row[0], 64)
				if err != nil {
					return "", "", nil, nil, fmt.Errorf("error parsing value: %v", err)
				}
				values = append(values, value)
			}
		} else {
			for _, row := range records {
				labels = append(labels, row[0])
				value, err := strconv.ParseFloat(row[1], 64)
				if err != nil {
					return "", "", nil, nil, fmt.Errorf("error parsing value: %v", err)
				}
				values = append(values, value)
			}
		}
	} else {
		if sepIndex != -1 {
			labels = args[:sepIndex]
			valuesArgs = args[sepIndex+1:]
			if len(labels) != len(valuesArgs) {
				return "", "", nil, nil, fmt.Errorf("number of labels and values must match. received: %d labels, %d values", len(labels), len(valuesArgs))
			}
		} else {
			valuesArgs = args
		}
	}

	for _, arg := range valuesArgs {
		value, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return "", "", nil, nil, fmt.Errorf("invalid number: %s\n%s", arg, usage)
		}
		values = append(values, value)
	}

	return title, chartType, labels, values, nil
}

func plot(values []float64, labels []string, width int, height int, chartType string) {
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

	barChartWidth := width - maxLabelLen - int(math.Log10(float64(maxValue+0.5))+1) - paddingWidth
	if barChartWidth < 1 {
		fmt.Println("warning: labels are too long to display chart values")
	}

	chartHeight := height - paddingHeight

	var characterUnit float64
	switch {
	case chartType == "bar":
		characterUnit = float64(maxValue) / float64(barChartWidth)
	case chartType == "column" || chartType == "col":
		characterUnit = float64(maxValue) / float64(chartHeight)
	default:
		characterUnit = float64(maxValue) / float64(barChartWidth)
	}
	var err error
	if chartType == "bar" {
		bar(values, labels, maxLabelLen, characterUnit)
	} else if chartType == "column" || chartType == "col" {
		err = column(values, labels, maxLabelLen, characterUnit, width)
	} else {
		bar(values, labels, maxLabelLen, characterUnit)
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
}

func getTerminalSize() (int, int, error) {
	cmd := exec.Command("tput", "cols")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	width, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, 0, err
	}

	cmd = exec.Command("tput", "lines")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err = cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func bar(values []float64, labels []string, maxLabelLen int, characterUnit float64) {
	for i, value := range values {
		if len(labels) > 0 {
			fmt.Printf("% *s│", maxLabelLen, labels[i])
		} else {
			fmt.Print("│")
		}

		scaledValue := int(float64(value)/characterUnit + 0.5)
		fmt.Print(blue)
		for j := 0; j < scaledValue; j++ {
			fmt.Print("■")
		}
		fmt.Print(reset)
		fmt.Printf(" %d\n", int(value+0.5))
	}
}

func column(values []float64, labels []string, maxLabelLen int, characterUnit float64, chartWidth int) (err error) {
	// check if too many columns, throw error
	if len(values)*2+1 > chartWidth {
		return fmt.Errorf("error plotting: too many columns to display, increase terminal width or try bar")
	}

	maxValue := getMaxValue(values)
	maxValueLen := len(strconv.Itoa(int(maxValue)))
	barWidth := min(max(maxLabelLen, maxValueLen, 2), 20)

	if (barWidth+1)*len(values)+1 > chartWidth {
		barWidth = chartWidth/len(values) - 1
	}
	chartHeight := int(maxValue/characterUnit + 0.5)
	for i := chartHeight + 1; i > 0; i-- {
		fmt.Print(strings.Repeat(" ", 1))
		for value := range values {
			if (values[value]/characterUnit + 0.5) >= float64(i) {
				fmt.Print(blue + strings.Repeat("█", barWidth) + reset)
			} else if (values[value]/characterUnit + 0.5) >= float64(i-1) {
				// check if value is longer than barWidth, if so, convert to scientific notation eg 1e+06
				if len(strconv.Itoa(int(values[value]))) > barWidth {
					if barWidth > 7 {
						fmt.Print(centerString(fmt.Sprintf("%.1e", values[value]), barWidth))
					} else if barWidth > 5 {
						fmt.Print(centerString(fmt.Sprintf("%.0e", values[value]), barWidth))
					} else {
						fmt.Print(centerString(" ", barWidth))
						err = fmt.Errorf("warning: Not enough space to display values")
					}
				} else {
					fmt.Print(centerString(strconv.Itoa(int(values[value])), barWidth))
				}

			} else {
				fmt.Print(strings.Repeat(" ", barWidth))
			}
			fmt.Print(strings.Repeat(" ", 1))
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", (barWidth+1)*len(values)+1))

	if len(labels) > 0 {
		fmt.Print(strings.Repeat(" ", 1))
		for _, label := range labels {
			if len(label) > barWidth {
				fmt.Print(label[:barWidth-1] + ".")
				err = fmt.Errorf("warning: labels truncated")
			} else {
				fmt.Print(centerString(label, barWidth))
			}
			fmt.Print(strings.Repeat(" ", 1))
		}
		fmt.Println()
	}
	if err != nil {
		return err
	}
	return nil
}

func getMaxValue(values []float64) float64 {
	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func centerString(s string, width int) string {
	padding := width - len(s)
	if padding > 0 {
		left := padding / 2
		right := padding - left
		return fmt.Sprintf("%*s%s%*s", left, "", s, right, "")
	}
	return s
}

func main() {
	title, chartType, labels, values, err := parseArgs(os.Args[1:])
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

	termWidth, termHeight, err := getTerminalSize()
	plotWidth := min(termWidth, maxPlotWidth)
	plotHeight := min(termHeight, maxPlotHeight)
	if err != nil {
		fmt.Printf("could not determine terminal size: %s\n", err)
		return
	}

	plot(values, labels, plotWidth, plotHeight, chartType)
}
