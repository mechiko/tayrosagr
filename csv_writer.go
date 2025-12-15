package main

import (
	"encoding/csv"
	"io"
	"strings"
)

type CustomCSVWriter struct {
	Writer    *csv.Writer
	Delimiter string
}

func NewCustomCSVWriter(w *csv.Writer, delimiter string) *CustomCSVWriter {
	return &CustomCSVWriter{
		Writer:    w,
		Delimiter: delimiter,
	}
}

func (c *CustomCSVWriter) Write(record []string) error {
	// Join the fields using the custom delimiter
	// joinedRecord := fmt.Sprintf("%s%s%s", record.code, c.Delimiter, record.serial)
	joinedRecord := strings.Join(record, c.Delimiter)

	// Write the joined record to the CSV writer
	return c.Writer.Write([]string{joinedRecord})
}

func NewWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = '\t'
	return
}

// func saveCsvCustom(name string, data [][]string) error {
// 	name = filepath.Join(outDir, name)
// 	file, err := os.Create(name)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	writer.Comma = '\t'
// 	writer.WriteAll(data) // calls Flush internally
// 	return writer.Error()
// }
