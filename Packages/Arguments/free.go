package Arguments

import (
	"SugarFree/Packages/Calculate"
	"SugarFree/Packages/Colors"
	"SugarFree/Packages/Reduce"
	"SugarFree/Packages/Utils"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// StageData represents data for each reduction stage
type StageData struct {
	stage   int
	entropy float64
}

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

		// Get variables from the command line
		file, _ := cmd.Flags().GetString("file")
		target, _ := cmd.Flags().GetFloat64("target")
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

		// Calculate initial entropy
		initialEntropy := Calculate.CalculateFullEntropy(originalData)

		// Display initial overall PE entropy
		fmt.Printf("[+] Initial Overall PE Entropy: %s\n", Colors.CalculateColor2Entropy(initialEntropy))

		// Create a slice to store entropy values for each stage
		stageData := []StageData{
			{0, initialEntropy}, // Include initial entropy as stage 0
		}

		// Copy the original data
		modifiedData := make([]byte, len(originalData))
		copy(modifiedData, originalData)

		// Add variables to track progress
		currentEntropy := initialEntropy
		iterationCount := 0
		lastEntropy := currentEntropy
		stuckCount := 0
		maxIterations := 10 // Maximum number of iterations to prevent infinite loops

		// Loop until we reach target entropy or can't reduce further
		for currentEntropy > target && iterationCount < maxIterations {
			// Call function named ApplyStrategy
			modifiedData = Reduce.ApplyStrategy(modifiedData, 60000)

			// Calculate new entropy
			currentEntropy = Calculate.CalculateFullEntropy(modifiedData)

			// Calculate reduction percentages
			totalReductionPercentage := ((initialEntropy - currentEntropy) / initialEntropy) * 100
			stageReductionPercentage := ((lastEntropy - currentEntropy) / lastEntropy) * 100

			// Calculate and display progress for this iteration
			iterationCount++

			// Store stage data for graphing
			stageData = append(stageData, StageData{
				stage:   iterationCount,
				entropy: currentEntropy,
			})

			// Convert current entropy to string for filename
			stageEntropy := strconv.FormatFloat(currentEntropy, 'f', 5, 64)

			// Build new filename for this stage
			stageFileName := Utils.BuildNewName(fileName, fileExtension, stageEntropy)

			// Write stage data to output file
			if err := ioutil.WriteFile(stageFileName, modifiedData, 0644); err != nil {
				fmt.Printf("[!] Error writing stage file: %v\n", err)
				continue
			}

			// Get file size for the new file
			newFileSize, err := Utils.GetFileSize(stageFileName)
			if err != nil {
				logger.Fatalf("[!] Error getting file size: %v\n", err)
				continue
			}

			if iterationCount == 1 {
				// For first stage, show entropy and current reduction percentage
				fmt.Printf("\n[+] Stage %d Reduction - Overall PE Entropy: %s\n",
					iterationCount,
					Colors.CalculateColor2Entropy(currentEntropy))
				fmt.Printf("[+] Stage %d Current Reduction Percentage: %s%%\n",
					iterationCount,
					Colors.BoldMagenta(fmt.Sprintf("%.2f", stageReductionPercentage)))
				fmt.Printf("[+] Stage %d File Size: %s KB\n",
					iterationCount,
					Colors.BoldYellow(newFileSize))
			} else {
				// For subsequent stages, show all messages in desired order
				fmt.Printf("\n[+] Stage %d Reduction - Overall PE Entropy: %s\n",
					iterationCount,
					Colors.CalculateColor2Entropy(currentEntropy))
				fmt.Printf("[+] Stage %d Current Reduction Percentage: %s%%\n",
					iterationCount,
					Colors.BoldMagenta(fmt.Sprintf("%.2f", stageReductionPercentage)))
				fmt.Printf("[+] Stage %d File Size: %s KB\n",
					iterationCount,
					Colors.BoldYellow(newFileSize))
				fmt.Printf("[+] Stage %d Total Reduction Percentage: %s%%\n",
					iterationCount,
					Colors.BoldBlue(fmt.Sprintf("%.2f", totalReductionPercentage)))
			}

			// Get absolute path for stage file
			stageFileName, err = Utils.GetAbsolutePath(stageFileName)
			if err != nil {
				logger.Fatalf("[!] Error getting absolute path for stage file: %v\n", err)
				continue
			}

			fmt.Printf("[+] Stage %d saved to: %s\n", iterationCount, Colors.BoldCyan(stageFileName))

			// Check if we're stuck (entropy isn't decreasing significantly)
			if lastEntropy-currentEntropy < 0.0001 {
				stuckCount++
				if stuckCount >= 3 { // If stuck for 3 iterations, break
					fmt.Printf("\n[!] Entropy reduction plateaued after %d stages\n", iterationCount)
					break
				}
			} else {
				stuckCount = 0
			}

			lastEntropy = currentEntropy
		}

		// If graph flag is enabled
		if graph {
			// Create a new plot
			p := plot.New()

			p.Title.Text = "Entropy Reduction Stages"
			p.X.Label.Text = "Stage"
			p.Y.Label.Text = "Entropy"

			// Create points for the line
			pts := make(plotter.XYs, len(stageData))
			for i, data := range stageData {
				pts[i].X = float64(data.stage)
				pts[i].Y = data.entropy
			}

			// Create a line plotter and set its style
			line, err := plotter.NewLine(pts)
			if err != nil {
				logger.Fatalf("\n[!] Error creating line plot: %v\n", err)
			} else {
				line.Color = color.RGBA{R: 255, B: 0, A: 255}
				line.Width = vg.Points(2)
				p.Add(line)

				// Add scatter points
				scatter, err := plotter.NewScatter(pts)
				if err != nil {
					logger.Fatalf("\n[!] Error creating scatter plot: %v\n", err)
				} else {
					scatter.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
					scatter.GlyphStyle.Radius = vg.Points(4)
					p.Add(scatter)

					getDateTime = time.Now().Format("20060102-150405")

					// Save the plot to a PNG file
					outputFile := fmt.Sprintf("%s_Entropy_Reduction_%s.png", fileName, getDateTime)
					if err := p.Save(8*vg.Inch, 6*vg.Inch, outputFile); err != nil {
						logger.Fatalf("\n[!] Error saving plot: %v\n", err)
					} else {
						fmt.Printf("\n[+] Entropy reduction graph saved to: %s\n", Colors.BoldYellow(outputFile))
					}
				}
			}
		}

		// Record the end time
		reductionEndTime := time.Now()

		// Calculate the duration
		reductionDurationTime := reductionEndTime.Sub(reductionStartTime)

		fmt.Printf("\n[*] Completed in: %s\n\n", Colors.BoldWhite(reductionDurationTime))

		return nil
	},
}
