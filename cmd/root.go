package cmd

import (
	"fmt"
	"os"

	"github.com/AdwaithSaiju/uopsy/internal/collector"
	"github.com/AdwaithSaiju/uopsy/internal/report"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("usage: uopsy [--json <path>]")
		os.Exit(1)
	}

	devices, err := collector.Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if os.Args[1] == "--json" {
		path := "report.json"
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		if err := report.WriteJSON(devices, path); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, d := range devices {
		fmt.Println(d.String())
	}
}
