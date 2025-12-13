package io

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func OutputResult(functionHeader []string, functionData [][]string, csvOutputFilePath string) error {
	if csvOutputFilePath != "" {
		if err := outputAsCSV(functionHeader, functionData, csvOutputFilePath); err != nil {
			return err
		}
	} else {
		if err := outputAsTable(functionHeader, functionData); err != nil {
			return err
		}
	}

	Logger.Info().Msgf("%d counts hit! ", len(functionData))

	return nil
}

func outputAsTable(header []string, data [][]string) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewTable(tableString,
		tablewriter.WithRendition(
			tw.Rendition{
				Symbols: tw.NewSymbols(tw.StyleASCII),
				Borders: tw.Border{
					Top:    tw.On,
					Bottom: tw.On,
					Left:   tw.On,
					Right:  tw.On,
				},
				Settings: tw.Settings{
					Separators: tw.Separators{
						BetweenRows: tw.On,
					},
					Lines: tw.Lines{
						ShowHeaderLine: tw.On,
					},
				},
			},
		),
	)

	table.Header(header)
	table.Bulk(data)
	table.Render()

	stringAsTableFormat := tableString.String()

	fmt.Fprintf(os.Stderr, "%s", stringAsTableFormat)

	return nil
}

func outputAsCSV(header []string, data [][]string, csvOutputFilePath string) error {
	file, err := os.Create(csvOutputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	var outputData [][]string

	outputData = append(outputData, header)
	outputData = append(outputData, data...)

	if err := w.WriteAll(outputData); err != nil {
		return err
	}

	if err := w.Error(); err != nil {
		return err
	}

	Logger.Info().Msg("Finished writing output!")

	return nil
}
