package Arguments

import (
	"SugarFree/Packages/Calculate"
	"SugarFree/Packages/Colors"
	"SugarFree/Packages/Reduce"
	"SugarFree/Packages/Utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// freeArgument represents the 'free' command in the CLI.
var freeArgument = &cobra.Command{
	// Use defines how the command should be called.
	Use:          "free",
	Short:        "Free command",
	Long:         "Lowers the overall entropy of a PE file",
	SilenceUsage: true,
	Aliases:      []string{"FREE", "Free"},

	// RunE defines the function to run when the command is executed.
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New(os.Stderr, "[!] ", 0)

		// Show ASCII banner
		ShowAscii()

		// Check if additional arguments were provided
		if len(os.Args) <= 2 {
			err := cmd.Help()
			if err != nil {
				logger.Fatal("Error ", err)
				return err
			}
			os.Exit(0)
		}

		// Define variables
		//var entropy string

		// Get variables from the command line
		file, _ := cmd.Flags().GetString("file")
		minimum, _ := cmd.Flags().GetFloat64("minimum")
		graph, _ := cmd.Flags().GetBool("graph")

		// Check if the file flag is empty
		if file == "" {
			logger.Fatal("Error: Input file is missing. Please provide it to continue...\n\n")
		}

		// Record start time for performance measurement
		reductionStartTime := time.Now()

		// Get the current date and time
		getDateTime := time.Now().Format("2006-01-02 15:04:05")

		fmt.Printf("[*] Starting PE entropy reduction on %s\n\n", Colors.BoldWhite(getDateTime))

		// Get absolute file path
		filePath, err := Utils.GetAbsolutePath(file)
		if err != nil {
			logger.Fatal("Error: ", err)
		}

		// Get file size
		fileSize, err := Utils.GetFileSize(filePath)
		if err != nil {
			logger.Fatal("Error: ", err)
		}

		// Get filename and extension
		fileName, fileExtension := Utils.SplitFileName(file)

		fmt.Printf("[+] Analyzing PE File: %s\n", Colors.BoldCyan(file))
		fmt.Printf("[+] Initial File Size: %s KB\n", Colors.BoldYellow(fileSize))

		// Read original binary data
		originalData, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}

		// Call function named CalculateFullEntropy
		initialEntropy := Calculate.CalculateFullEntropy(originalData)

		// Display initial overall PE entropy
		fmt.Printf("[+] Initial Overall PE Entropy: %s\n\n", Colors.CalculateColor2Entropy(initialEntropy))

		// Copy the original data
		modifiedData := make([]byte, len(originalData))
		copy(modifiedData, originalData)

		// Call function named ApplyStrategy
		modifiedData = Reduce.ApplyStrategy(modifiedData, 60000)

		finalEntropy := Calculate.CalculateFullEntropy(modifiedData)
		fmt.Printf("[+] Staged Reduction - Overall PE Entropy: %s\n", Colors.CalculateColor2Entropy(finalEntropy))
		reductionPercentage := ((initialEntropy - finalEntropy) / initialEntropy) * 100
		fmt.Printf("[+] Entropy Reduction Percentage: %s%%\n", Colors.BoldBlue(fmt.Sprintf("%.2f", reductionPercentage)))

		// Convert float to string
		entropy := strconv.FormatFloat(finalEntropy, 'f', -1, 64)

		// Call function named BuildNewName
		newFileName := Utils.BuildNewName(fileName, fileExtension, entropy)

		// Write modified data to output file
		if err := ioutil.WriteFile(newFileName, modifiedData, 0644); err != nil {
			fmt.Printf("[!] Error writing output file: %v\n", err)
			os.Exit(1)
		}

		// Call function named GetAbsoluteFilePath
		newFileName, err = Utils.GetAbsolutePath(newFileName)
		if err != nil {
			logger.Fatal("Error: ", err)
		}

		// Display the new file name
		fmt.Printf("[+] Saved to: %s\n\n", Colors.BoldCyan(newFileName))

		// Record the end time
		reductionEndTime := time.Now()

		// Calculate the duration
		reductionDurationTime := reductionEndTime.Sub(reductionStartTime)

		fmt.Printf("[*] Completed in: %s\n\n", Colors.BoldWhite(reductionDurationTime))

		//////// To be removed //////////
		fmt.Print(minimum, graph)

		return nil
	},
}
