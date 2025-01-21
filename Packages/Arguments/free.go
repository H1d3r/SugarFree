package Arguments

import (
	"SugarFree/Packages/Colors"
	"SugarFree/Packages/Entropy"
	"SugarFree/Packages/Graph"
	"SugarFree/Packages/Manager"
	"SugarFree/Packages/Output"
	"SugarFree/Packages/Utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// freeArgument represents the 'free' command in the CLI.
var freeArgument = &cobra.Command{
	// Use defines how the command should be called.
	Use:          "free",
	Short:        "Reduce and analyze file entropy",
	Long:         "Analyze and optionally reduce entropy in PE files while maintaining functionality",
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

		// Get command line flags
		file, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")
		graph, _ := cmd.Flags().GetBool("graph")

		// Validate input file
		if file == "" {
			logger.Fatal("Error: Input file is missing. Please provide it to continue...\n")
		}

		// Record start time for performance measurement
		startTime := time.Now()

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

		// Read PE sections and calculate initial entropy
		fmt.Printf("[+] Analyzing the PE file: %s\n", Colors.BoldGreen(file))
		fmt.Printf("[+] Initial file size: %s bytes\n", Colors.BoldYellow(fileSize))

		sections, err := Manager.ReadSections(filePath)
		if err != nil {
			logger.Fatal("Error: ", err)
		}

		// Calculate and display initial entropy
		initialEntropy := Manager.FullEntropy(sections)
		fmt.Printf("[+] Initial PE entropy: %s\n", Colors.BoldCyan(fmt.Sprintf("%.5f", initialEntropy)))

		// Process sections for entropy reduction
		reducedSections, err := Entropy.ReduceEntropy(filePath, Entropy.ReductionStrategy{Aggressive: true})
		if err != nil {
			logger.Fatal("Error during entropy reduction: ", err)
		}

		// Calculate and display final entropy
		finalEntropy := Manager.FullEntropy(reducedSections)
		fmt.Printf("[+] Final PE entropy: %s\n", Colors.BoldGreen(fmt.Sprintf("%.5f", finalEntropy)))

		// Display entropy reduction percentage
		reductionPercentage := ((initialEntropy - finalEntropy) / initialEntropy) * 100
		fmt.Printf("[+] Entropy reduction: %s%%\n", Colors.BoldYellow(fmt.Sprintf("%.2f", reductionPercentage)))

		// Generate entropy visualization if requested
		if graph {
			fmt.Print("[+] Generating entropy visualization graph...\n")
			graphPath := file + "_entropy_graph.png"
			err = Graph.GenerateEntropyGraph(sections, reducedSections, graphPath)
			if err != nil {
				logger.Printf("Warning: Failed to generate entropy graph: %v\n", err)
			} else {
				fmt.Printf("[+] Entropy graph saved to: %s\n", Colors.BoldGreen(graphPath))
			}
		}

		// Save results to output file if specified
		if output != "" {
			outputPath, err := Utils.GetAbsolutePath(output)
			if err != nil {
				logger.Fatal("Error: ", err)
			}

			// Convert to Output.Section format
			var outputSections []Output.Section
			for _, section := range reducedSections {
				outputSections = append(outputSections, Output.Section{
					Name:    section.Name,
					Entropy: section.Entropy,
				})
			}

			// Write results to file
			Output.WriteEntropyReport(outputSections, outputPath, file, fileSize, initialEntropy, finalEntropy)
			fmt.Printf("\n[+] Results saved to: %s\n", Colors.BoldGreen(outputPath))
		}

		// Display execution time
		duration := time.Since(startTime)
		fmt.Printf("\n[+] Completed in: %s\n\n", Colors.BoldWhite(duration))

		return nil
	},
}
