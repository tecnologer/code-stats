# Code Stats

Collects the code statistics of a given directory.

## Pre-requisites

- `scc` in the PATH, to install it run the following command:
  ```shell
  go install github.com/boyter/scc/v3@latest
  ```
  
## Usage

```text
Usage of code-stats:
  -draw-chart
        Draw chart
  -input string
        Path to the input file, separated by commas, could be used with stdin (default ".stats")
  -languages string
        Languages to include in the chart, require at least one and --draw-chart (default "go")
  -omit-dir string
        Directories to omit from the stats (default ".idea,vendor,.stats")
  -only-compare-input
        Only compare the input files, do not calculate the current stats
  -version
        Show version
```
