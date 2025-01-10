package Arguments

import (
	"SugarFree/Packages/Colors"
	"SugarFree/Packages/Manager"
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
	Short:        "info command",
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

		// Get variables from the command line.
		file, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")

		// Check if the file flag is empty.
		if file == "" {
			logger.Fatal("Error: Input file is missing. Please provide it to continue...\n")
		}

		// Record the start time
		calculateStartTime := time.Now()

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
		sections, err := Manager.ReadSections(filePath)
		if err != nil {
			log.Fatal(err)
		}

		// Convert Manager.SectionEntropy to Output.Section
		var outputSections []Output.Section
		for _, section := range sections {
			outputSections = append(outputSections, Output.Section{
				Name:    section.Name,
				Entropy: section.Entropy,
			})
		}

		// Call function named FullEntropy
		fullEntropy := Manager.FullEntropy(sections)

		fmt.Printf("[+] Analyzing the PE file: %s\n", Colors.BoldGreen(file))
		fmt.Printf("[+] File Size: %s bytes\n", Colors.BoldYellow(fileSize))
		fmt.Printf("[+] Full PE Entropy: %s\n", Colors.BoldCyan(fmt.Sprintf("%.5f", fullEntropy)))
		fmt.Print("[+] PE Sections and their entropy:\n")
		for _, section := range sections {
			// Call function ColorManager
			sectionName := Manager.ColorNameManager(section.Name)

			// Check if the entropy is less than 5.0
			if section.Entropy < 5.0 {
				sectionEntropy = Colors.BoldGreen(fmt.Sprintf("%.5f", section.Entropy))
			} else {
				sectionEntropy = Colors.BoldRed(fmt.Sprintf("%.5f", section.Entropy))
			}

			fmt.Printf("	>>> \"%s\" Scored Entropy Of Value: %s\n", sectionName, sectionEntropy)
		}

		// Check if the output flag is empty.
		if output != "" {
			// Call function named WriteToFile
			Output.WriteToFile(outputSections, output, file, fileSize)

			// Call function named GetAbsolutePath
			outputFilePath, err := Utils.GetAbsolutePath(output)
			if err != nil {
				logger.Fatal("Error: ", err)
			}

			fmt.Printf("\n[+] Results saved to: %s\n", Colors.BoldGreen(outputFilePath))
		}

		// Record the end time
		calculateEndTime := time.Now()

		// Calculate the duration
		calculateDurationTime := calculateEndTime.Sub(calculateStartTime)

		// Print the duration
		fmt.Printf("\n[+] Completed in: %s\n\n", Colors.BoldWhite(calculateDurationTime))

		return nil
	},
}
