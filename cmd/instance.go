package cmd

import (
	"github.com/spf13/cobra"
)

var instanceTypeFamilies = []string{
	"cpu",
	"gpu",
	"gpu2",
	"memory",
	"standard",
	"storage",
}

var instanceTypeSizes = []string{
	"micro",
	"tiny",
	"small",
	"medium",
	"large",
	"extra-large",
	"huge",
	"jumbo",
	"mega",
	"titan",
}

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Compute instances management",
}

func init() {
	computeCmd.AddCommand(instanceCmd)
}
