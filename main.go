package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mtnmunuklu/bridge/sigma"
	"github.com/mtnmunuklu/bridge/sigma/evaluator"
)

var (
	filePath      string
	configPath    string
	fileContent   string
	configContent string
	showHelp      bool
	outputJSON    bool
	outputPath    string
	version       bool
	caseSensitive bool
)

func printUsage() {
	fmt.Println("Usage: bridge -filepath <path> -config <path> [flags]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("Example:")
	fmt.Println("  bridge -filepath /path/to/file -config /path/to/config")
}

// Set up the command-line flags
func init() {
	flag.StringVar(&filePath, "filepath", "", "Name or path of the file or directory to read")
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.StringVar(&fileContent, "filecontent", "", "Base64-encoded content of the file or directory to read")
	flag.StringVar(&configContent, "configcontent", "", "Base64-encoded content of the configuration file")
	flag.BoolVar(&showHelp, "help", false, "Show usage")
	flag.BoolVar(&outputJSON, "json", false, "Output results in JSON format")
	flag.StringVar(&outputPath, "output", "", "Output directory for writing files")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.BoolVar(&caseSensitive, "cs", false, "Case sensitive mode")
	flag.Parse()

	// If the version flag is provided, print version information and exit
	if version {
		fmt.Println("Bridge version 1.0.0")
		os.Exit(1)
	}

	// If the help flag is provided, print usage information and exit
	if showHelp {
		printUsage()
		os.Exit(1)
	}

	// Check if filepath and configpath are provided as command-line arguments
	if flag.NArg() > 0 {
		filePath = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		configPath = flag.Arg(1)
	}

	// Check if both filecontent and configcontent are provided
	if (filePath == "" && fileContent == "") || (configPath == "" && configContent == "") {
		fmt.Println("Please provide either file paths or file contents, and either config path or config content.")
		printUsage()
		os.Exit(1)
	}
}

func formatSigmaJSONResult(rule sigma.Rule, queries map[int]string) []byte {
	// Define a struct type named JSONResult to represent the JSON output fields.
	type JSONResult struct {
		Name           string   `json:"Name"`
		Description    string   `json:"Description"`
		Query          string   `json:"Query"`
		InsertDate     string   `json:"InsertDate"`
		LastUpdateDate string   `json:"LastUpdateDate"`
		Tags           []string `json:"Tags"`
		Level          string   `json:"Level"`
	}

	// Create a strings.Builder variable named query.
	var query strings.Builder
	for i, qry := range queries {
		// Add a newline character if the index is greater than zero.
		if i > 0 {
			query.WriteString("\n")
		}
		query.WriteString(qry)
	}

	// Create an instance of the JSONResult struct.
	jsonResult := JSONResult{
		Name:           rule.Title,
		Description:    rule.Description + "\n\nAuthor: " + rule.Author + "\nSigma Repository: [GitHub](https://github.com/SigmaHQ/sigma)",
		Query:          query.String(),
		InsertDate:     time.Now().UTC().Format(time.RFC3339),
		LastUpdateDate: time.Now().UTC().Format(time.RFC3339),
		Tags:           rule.Tags,
		Level:          rule.Level,
	}

	// Marshal the JSONResult struct into JSON data.
	jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil
	}

	return jsonData
}

func main() {

	// Read the contents of the file(s) specified by the filepath flag or filecontent flag
	fileContents := make(map[string][]byte)
	var err error

	// Check if file paths are provided
	if filePath != "" {
		// Check if the filepath is a directory
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error getting file/directory info:", err)
			return
		}

		if fileInfo.IsDir() {
			// FilePath is a directory, so walk the directory to read all the files inside it
			filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Println("Error accessing file:", err)
					return nil
				}
				if !info.IsDir() {
					// read file content
					content, err := os.ReadFile(path)
					if err != nil {
						fmt.Println("Error reading file:", err)
						return nil
					}
					fileContents[path] = content
				}
				return nil
			})
		} else {
			// FilePath is a file, so read its contents
			fileContents[filePath], err = os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
		}
	} else if fileContent != "" {
		// Check if the filecontent is a directory
		lines := strings.Split(fileContent, "\n")
		if len(lines) > 1 {
			// Read all lines as separate files
			for _, line := range lines {
				// Decode base64 content
				decodedContent, err := base64.StdEncoding.DecodeString(line)
				if err != nil {
					fmt.Println("Error decoding base64 content:", err)
					return
				}
				fileContents[line] = decodedContent
			}
		} else {
			// Read its content
			// Decode base64 content
			decodedContent, err := base64.StdEncoding.DecodeString(fileContent)
			if err != nil {
				fmt.Println("Error decoding base64 content:", err)
				return
			}
			fileContents["filecontent"] = decodedContent
		}
	}

	// Read the contents of the configuration file or use configcontent
	var configContents []byte
	if configPath != "" {
		configContents, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading configuration file:", err)
			return
		}
	} else if configContent != "" {
		// decode base64 content
		decodedContent, err := base64.StdEncoding.DecodeString(configContent)
		if err != nil {
			fmt.Println("Error decoding base64 content:", err)
			return
		}
		configContents = decodedContent
	}

	for _, fileContent := range fileContents {
		sigmaRule, err := sigma.ParseRule(fileContent)
		if err != nil {
			fmt.Println("Error parsing rule:", err)
			continue
		}

		// Parse the configuration file as a Sigma config
		config, err := sigma.ParseConfig(configContents)
		if err != nil {
			fmt.Println("Error parsing config:", err)
			continue
		}

		var rule *evaluator.RuleEvaluator

		if caseSensitive {
			// Evaluate the Sigma rule against the config using case sensitive mode
			rule = evaluator.ForRule(sigmaRule, evaluator.WithConfig(config), evaluator.CaseSensitive)
		} else {
			// Evaluate the Sigma rule against the config
			rule = evaluator.ForRule(sigmaRule, evaluator.WithConfig(config))
		}

		result, err := rule.Bridges()
		if err != nil {
			fmt.Println("Error converting rule:", err)
			continue
		}

		var output string

		// Print the results of the query
		if outputJSON {
			jsonResult := formatSigmaJSONResult(sigmaRule, result.QueryResults)
			output = string(jsonResult)
		} else {
			var builder strings.Builder
			for _, queryResult := range result.QueryResults {
				builder.WriteString(queryResult + "\n")
			}
			output = builder.String()
		}

		// Check if outputPath is provided
		if outputPath != "" {
			// Create the output file path using the Name field from the rule
			outputFilePath := filepath.Join(outputPath, fmt.Sprintf("%s.json", sigmaRule.Title))

			// Write the output string to the output file
			err := os.WriteFile(outputFilePath, []byte(output), 0644)
			if err != nil {
				fmt.Println("Error writing output to file:", err)
				continue
			}

			fmt.Printf("Output for rule '%s' written to file: %s\n", sigmaRule.Title, outputFilePath)
		} else {
			fmt.Printf("%s\n", sigmaRule.Title)
			fmt.Printf("%s", output)
		}

	}

}
