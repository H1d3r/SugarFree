package Arguments

import (
	"SugarFree/Packages/Colors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	__version__      = "2.0"
	__license__      = "MIT"
	__author__       = []string{"@nickvourd", "@IAMCOMPROMISED"}
	__github__       = "https://github.com/nickvourd/SugarFree"
	__version_name__ = "Ocean Words"
	__ascii__        = `
 
███████╗██╗   ██╗ ██████╗  █████╗ ██████╗ ███████╗██████╗ ███████╗███████╗
██╔════╝██║   ██║██╔════╝ ██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔════╝██╔════╝
███████╗██║   ██║██║  ███╗███████║██████╔╝█████╗  ██████╔╝█████╗  █████╗  
╚════██║██║   ██║██║   ██║██╔══██║██╔══██╗██╔══╝  ██╔══██╗██╔══╝  ██╔══╝  
███████║╚██████╔╝╚██████╔╝██║  ██║██║  ██║██║     ██║  ██║███████╗███████╗
╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝╚══════╝╚══════╝                                                        
`

	__text__ = `
SugarFree v%s - Less sugar (entropy) for your binaries.
SugarFree is an open source tool licensed under %s.
Written with <3 by %s && %s...
Please visit %s for more...

`

	SugarFreeCli = &cobra.Command{
		Use:          "SugarFree",
		SilenceUsage: true,
		RunE:         StartSugarFree,
		Aliases:      []string{"sugarfree", "SUGARFREE", "sf"},
	}
)

// ShowAscii function
func ShowAscii() {
	// Initialize RandomColor
	randomColor := Colors.RandomColor()
	fmt.Print(randomColor(__ascii__))
	fmt.Printf(__text__, __version__, __license__, __author__[0], __author__[1], __github__)
}

// init function
// init all flags.
func init() {
	// Disable default command completion for SugarFree CLI.
	SugarFreeCli.CompletionOptions.DisableDefaultCmd = true

	// Add commands to the SugarFree CLI.
	SugarFreeCli.Flags().SortFlags = true
	SugarFreeCli.Flags().BoolP("version", "v", false, "Show SugarFree current version")
	SugarFreeCli.AddCommand(infoArgument)
	SugarFreeCli.AddCommand(freeArgument)

	// Add flags to the 'info' command.
	infoArgument.Flags().SortFlags = true
	infoArgument.Flags().StringP("file", "f", "", "Set input file")
	infoArgument.Flags().StringP("output", "o", "", "Save results to output file")

	// Add flags to the 'free' command.
	freeArgument.Flags().SortFlags = true
	freeArgument.Flags().StringP("file", "f", "", "Set input file")
	freeArgument.Flags().Float64P("target", "t", 4.6, "Set target entropy value to achieve")
	freeArgument.Flags().StringP("strategy", "s", "zero", "Set strategy to apply (i.e., zero, word)")
	freeArgument.Flags().BoolP("graph", "g", false, "Enable entropy graph")
}

// ShowVersion function
func ShowVersion(versionFlag bool) {
	// if one argument
	if versionFlag {
		// if version flag exists
		fmt.Print("[+] Current version: " + Colors.BoldRed(__version__) + "\n\n[+] Version name: " + Colors.BoldRed(__version_name__) + "\n\n")
		os.Exit(0)
	}
}

// StartSugarFree function
func StartSugarFree(cmd *cobra.Command, args []string) error {
	logger := log.New(os.Stderr, "[!] ", 0)

	// Call function named ShowAscii
	ShowAscii()

	// Check if additional arguments were provided.
	if len(os.Args) == 1 {
		// Display help message.
		err := cmd.Help()

		// If error exists
		if err != nil {
			logger.Fatal("Error: ", err)
			return err
		}
	}

	// Obtain flag
	versionFlag, _ := cmd.Flags().GetBool("version")

	// Call function named ShowVersion
	ShowVersion(versionFlag)

	return nil
}
