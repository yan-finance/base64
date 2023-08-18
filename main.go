package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  base64 [OPTIONS] [TEXT]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	decodeFlag := flag.Bool("d", false, "Decode incoming Base64 data.")
	inputFile := flag.String("i", "", "Input file")
	outputFile := flag.String("o", "", "Output file (default: stdout)")
	breakCount := flag.Int("b", 0, "Break encoded string into num character lines")
	helpFlag := flag.Bool("h", false, "Print usage summary and exit.") // New help flag

	flag.Parse()

	fmt.Printf("%v", len(os.Args))
	if len(os.Args) <= 1 || *helpFlag {
		printUsage()
		return
	}

	var input io.Reader
	var output io.Writer = os.Stdout // Default to stdout if not specified

	defer func() {
		fmt.Println()
	}()

	if *decodeFlag {

		if *inputFile == "" || *inputFile == "-" {
			// Use command line arguments as input text for decoding
			inputText := strings.Join(flag.Args(), " ")
			inputText = strings.TrimSpace(inputText) // Trim leading/trailing spaces
			if inputText != "" {
				input = strings.NewReader(inputText)
				decoder := base64.NewDecoder(base64.StdEncoding, input)
				_, err := io.Copy(output, decoder)
				if err != nil {
					fmt.Println("Error decoding input:", err)
					return
				}
			}
		} else {
			file, err := os.Open(*inputFile)
			if err != nil {
				fmt.Println("Error opening input file:", err)
				return
			}
			defer file.Close()
			input = file
			decoder := base64.NewDecoder(base64.StdEncoding, input)
			_, err = io.Copy(output, decoder)
			fmt.Println()
			if err != nil {
				fmt.Println("Error decoding input:", err)
				return
			}
		}

	} else {
		inputText := strings.Join(flag.Args(), " ") // Combine remaining args as input text

		if *inputFile == "" || *inputFile == "-" {
			input = strings.NewReader(inputText)
		} else {
			file, err := os.Open(*inputFile)
			if err != nil {
				fmt.Println("Error opening input file:", err)
				return
			}
			defer file.Close()
			input = file
		}

		if *outputFile != "" {
			file, err := os.Create(*outputFile)
			if err != nil {
				fmt.Println("Error creating output file:", err)
				return
			}
			defer file.Close()
			output = file
		}

		if *breakCount > 0 {
			lineBreakWriter := NewLineBreakWriter(output, *breakCount)
			encoder := base64.NewEncoder(base64.StdEncoding, lineBreakWriter)
			defer encoder.Close()
			_, err := io.Copy(encoder, input)
			if err != nil {
				fmt.Println("Error encoding input:", err)
				return
			}
		} else {
			encoder := base64.NewEncoder(base64.StdEncoding, output)
			defer encoder.Close()
			_, err := io.Copy(encoder, input)
			if err != nil {
				fmt.Println("Error encoding input:", err)
				return
			}
		}
	}
}

type LineBreakWriter struct {
	writer    io.Writer
	breakFreq int
	count     int
}

func NewLineBreakWriter(writer io.Writer, breakFreq int) *LineBreakWriter {
	return &LineBreakWriter{
		writer:    writer,
		breakFreq: breakFreq,
	}
}

func (w *LineBreakWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		n, err = w.writer.Write([]byte{b})
		w.count++
		if w.count == w.breakFreq {
			w.count = 0
			w.writer.Write([]byte{'\n'})
		}
	}
	return n, err
}
