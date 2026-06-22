package cmd

import (
	"fmt"
	"os"
	"strings"
)

const asciiArt = `
+===========================================================+
|                                                           |
|        ██╗   ██╗ ██████╗ ██████╗ ███████╗██╗   ██╗        |
|        ██║   ██║██╔═══██╗██╔══██╗██╔════╝╚██╗ ██╔╝        |
|        ██║   ██║██║   ██║██████╔╝███████╗ ╚████╔╝         |
|        ██║   ██║██║   ██║██╔═══╝ ╚════██║  ╚██╔╝          |
|        ╚██████╔╝╚██████╔╝██║     ███████║   ██║           |
|         ╚═════╝  ╚═════╝ ╚═╝     ╚══════╝   ╚═╝           |
|                                                           |
|              USB Forensics & Analysis Tool                |
+===========================================================+
`

func Execute() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "scan", "-s", "--scan":
		Scan()
	case "-h", "--help", "help":
		showHelp()
	default:
		fmt.Printf("Unknown command: '%s'\n\n", os.Args[1])
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Print(asciiArt)
	fmt.Println()
	fmt.Println("  USAGE:")
	fmt.Println("    uopsy <command>")
	fmt.Println()
	fmt.Println("  COMMANDS:")
	fmt.Println("    scan     Scan and list connected USB devices")
	fmt.Println("    help     Show this help message")
	fmt.Println()
	fmt.Println("  EXAMPLE:")
	fmt.Println("    uopsy scan")
	fmt.Println()
	fmt.Println("  For more information, visit: https://github.com/AdwaithSaiju/uopsy")
}
