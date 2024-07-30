# Code Stats

Collects the code statistics of a given directory, and could compare with the previous stats.

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
        Languages to include in the chart, require at least one and --draw-chart
  -omit-dir string
        Directories to omit from the stats (default ".idea,vendor,.stats")
  -only-compare-input
        Only compare the input files, do not calculate the current stats
  -stat-name string
        Name of the stat, accepted values: bytes, code_bytes, lines, code, comment, blank, complexity, count, weighted_complexity (default "code")
  -version
        Show version
```
