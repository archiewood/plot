# Plot

A minimalist CLI tool to plot data to the terminal.

- Data: Plot data from CSV files or command line arguments.
- Charts: Currently only support bar charts.

## Usage

```bash
plot [-t title] file.csv
plot [-t title] [labels ,]  values
```

### Bar Chart from CSV

```bash
plot test/two_col_header.csv
```

### Bar Chart with Title and Labels

```bash
plot -t "Programming Languages" go python r c++ , 84 950 923 27
```


## Credits

This tool is heavily inspired by [YouPlot](https://github.com/red-data-tools/YouPlot)