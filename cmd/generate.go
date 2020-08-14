package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/luncj/mess/generator"
	"github.com/luncj/mess/schema"
)

func generateCmd() *cobra.Command {

	var dml string
	var numRows uint
	var schemaPath string
	var metadataPath string
	var outputPath string

	cmd := cobra.Command{
		Use:   "generate",
		Short: "Generate SQL with random data",
		Run: func(*cobra.Command, []string) {
			output, err := os.Create(outputPath)
			if err != nil {
				log.Fatalf("failed to open output file: %s", err)
			}

			defer output.Close()

			s, err := schema.FromFile(schemaPath)
			if err != nil {
				log.Fatalf("failed to read schema from file: %s", err)
			}

			d, err := schema.DMLFromString(dml)
			if err != nil {
				log.Fatalf("failed to fetch DML: %s", err)
			}

			var md *generator.Metadata
			md, err = generator.Open(s, metadataPath)
			if err != nil {
				log.Fatalf("failed to read metadata: %s", err)
			}
			defer md.Close()

			g := generator.NewMySQLGenerator()
			sqls, err := g.Generate(md, s, d, numRows)
			_, err = output.WriteString(strings.Join(sqls, "\n"))
			if err != nil {
				log.Fatalf("failed to write SQL to output: %s", err)
			}
		},
	}

	f := cmd.PersistentFlags()
	f.StringVar(&dml, "dml", "update", "DML")
	f.StringVar(&schemaPath, "schema-path", "./examples/insert/schema.json", "Path of schema definition")
	f.StringVar(&metadataPath, "metadata-path", "./metadata.json", "Path of metadata for increment generating")
	f.StringVar(&outputPath, "output-path", "output.sql", "Path of generated SQL")
	f.UintVar(&numRows, "num-rows", 1_000, "Number of rows")

	return &cmd
}

func init() {
	rootCmd.AddCommand(generateCmd())
}
