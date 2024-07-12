package cmd

import (
	"bufio"
	"compress/gzip"
	"os"

	"github.com/cuducos/minha-receita/transform"
	"github.com/spf13/cobra"
)

const exportHelper = `
Convert the CSV files from the Federal Revenue for venues (ESTABELE group of
files) into one big JSON file, joining information from all other source CSV files.`

const exportPath = "export.json.gz"

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Transforms the CSV files into JSON file",
	Long:  exportHelper,
	RunE: func(_ *cobra.Command, _ []string) error {
		if err := assertDirExists(); err != nil {
			return err
		}
		fn, err := os.Create(exportPath)
		if err != nil {
			return err
		}
		defer fn.Close()
		gz := gzip.NewWriter(fn)
		defer gz.Close()
		out := bufio.NewWriter(gz)
		defer out.Flush()
		return transform.Export(dir, out, batchSize, highMemory)
	},
}

func exportCLI() *cobra.Command {
	exportCmd = addDataDir(exportCmd)
	exportCmd.Flags().StringVarP(
		&targetDir,
		"output",
		"o",
		exportPath,
		"directory for the output file",
	)
	exportCmd.Flags().IntVarP(&batchSize, "batch-size", "b", transform.BatchSize, "size of the batch to save to the database")
	exportCmd.Flags().BoolVarP(&highMemory, "high-memory", "x", highMemory, "high memory availability mode, faster but requires a lot of free RAM")
	return exportCmd
}
