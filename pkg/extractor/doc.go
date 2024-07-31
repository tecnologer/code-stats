// Package extractor provides functionality to extract and process statistical
// data from various sources including files and standard input. It supports
// reading JSON formatted data, merging data from multiple sources, and
// interfacing with external tools for data generation.
//
// The package is designed to be used with the code-stats system, facilitating
// the aggregation and analysis of code metrics across multiple files or
// repositories. It includes functions to handle input from both file paths
// and standard input, allowing for flexible data integration.
//
// This package utilizes the 'models' package for structuring the statistical data
// and the 'file' package for file operations, ensuring a cohesive and modular
// approach to data handling within the code-stats system.
package extractor
