# Plot

A CLI tool to plot data to the terminal.

Currently support bar charts.

## Usage

```bash
plot [-t "title"] [labels,] data
```

### Bar Chart with Title and Labels

```bash
plot -t "Programming Languages" go python r c++ , 84 950 923 27
```

### Bar Chart from CSV

```bash
plot test/two_col_header.csv
```

## Credits

This tool is heavily inspired by [YouPlot](https://github.com/red-data-tools/YouPlot)