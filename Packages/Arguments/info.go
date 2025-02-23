package Arguments

import (
	"SugarFree/Packages/Calculate"
	"SugarFree/Packages/Colors"
	"SugarFree/Packages/Output"
	"SugarFree/Packages/Utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// infoArgument represents the 'info' command in the CLI.
var infoArgument = &cobra.Command{
	// Use defines how the command should be called.
	Use:          "info",
	Short:        "Info command",
	Long:         "Calculates the entropy of a PE file and its sections",
	SilenceUsage: true,
	Aliases:      []string{"INFO", "Info"},

	// RunE defines the function to run when the command is executed.
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New(os.Stderr, "[!] ", 0)

		// Call function named ShowAscii
		ShowAscii()

		// Check if additional arguments were provided.
		if len(os.Args) <= 2 {
			// Show help message.
			err := cmd.Help()
			if err != nil {
				logger.Fatal("Error ", err)
				return err
			}

			// Exit the program.
			os.Exit(0)
		}

		// Define variables
		var sectionEntropy string

		// Get variables from the command line
		file, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")

		// Check if the file flag is empty
		if file == "" {
			logger.Fatal("Error: Input file is missing. Please provide it to continue...\n\n")
		}

		// Record the start time
		calculateStartTime := time.Now()

		// Get the current date and time
		getDateTime := time.Now().Format("2006-01-02 15:04:05")

		fmt.Printf("[*] Starting PE analysis on %s\n\n", Colors.BoldWhite(getDateTime))

		// Call function named GetAbsolutePath
		filePath, err := Utils.GetAbsolutePath(file)
		if err != nil {
			logger.Fatal("Error: ", err)
		}

		// Call function named GetFileSize
		fileSize, err := Utils.GetFileSize(filePath)
		if err != nil {
			logger.Fatal("Error: ", err)
		}

		// Call function named ReadSections
		sections, err := Calculate.ReadSections(filePath)
		if err != nil {
			log.Fatal(err)
		}

		// Convert Calcualte.SectionEntropy to Output.Section
		var outputSections []Output.Section
		for _, section := range sections {
			outputSections = append(outputSections, Output.Section{
				Name:    section.Name,
				Entropy: section.Entropy,
			})
		}

		// Call function named FullEntropy
		fullEntropy := Calculate.CalculateFullEntropy(sections)

		// Print the results
		fmt.Printf("[+] Analyzing PE File: %s\n", Colors.BoldCyan(file))
		fmt.Printf("[+] File Size: %s KB\n", Colors.BoldYellow(fileSize))
		fmt.Printf("[+] Overall PE Entropy: %s\n\n", Colors.CalculateColor2Entropy(fullEntropy))
		fmt.Print("[+] PE Sections Entropy:\n")
		for _, section := range sections {
			// Call function ColorManager
			sectionName := Colors.ColorNameManager(section.Name)

			// Call function named CalculateColor2Entropy
			sectionEntropy = Colors.CalculateColor2Entropy(section.Entropy)

			// Print the results
			fmt.Printf("	>>> \"%s\" Scored Entropy Of Value: %s\n", sectionName, sectionEntropy)
		}

		// Check if the output flag is empty.
		if output != "" {
			// Call function named WriteToFile
			Output.Write2File(outputSections, output, file, fileSize, fullEntropy, getDateTime)

			// Call function named GetAbsolutePath
			outputFilePath, err := Utils.GetAbsolutePath(output)
			if err != nil {
				logger.Fatal("Error: ", err)
			}

			fmt.Printf("\n[+] Results saved to: %s\n\n", Colors.BoldCyan(outputFilePath))
		}

		// Record the end time
		calculateEndTime := time.Now()

		// Calculate the duration
		calculateDurationTime := calculateEndTime.Sub(calculateStartTime)

		// Print the duration
		fmt.Printf("[*] Completed in: %s\n\n", Colors.BoldWhite(calculateDurationTime))

		return nil
	},
}
