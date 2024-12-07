package main

import (
	"fmt"
	"os"
	"strings"
	"encoding/csv"
	"time"
)

import flag "github.com/spf13/pflag"

func makeCSV (fileName string, force bool) {
	if !strings.HasSuffix(fileName, ".csv") {
		fmt.Println("File does not ends with .csv. Appending .csv to the file name")
		fileName = fileName + ".csv"
	}

	_, existsErr := os.Stat(fileName)
	if existsErr == nil {
		if !force {
			fmt.Println("File already exists. Use --make-force or -M to overwrite the file")
			return
		} else {
			fmt.Println("File already exists. Overwriting the file")
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	writer := csv.NewWriter(file)
	writer.Write([]string{"Task", "Time"})
	writer.Flush()
	defer file.Close()
	fmt.Println("File created successfully")
}

func displayFile(fileName string) {
	fmt.Println("Displaying file: ", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		fmt.Println("File is empty")
		return
	}
	for _, row := range data {
		for _, col := range row {
			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}
}

func writeLine(fileName string, text string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	data := []string{text, time.Now().Format("2006-01-02 15:04:05")}
	writer := csv.NewWriter(file)
	writer.Write(data)
	writer.Flush()
}

func main() {
	var makeFileName = flag.StringP("make", "m", "", "The name of file to make")
	var makeForce = flag.BoolP("make-force", "M", false, "Force overwrite the file")
	var displayFileName = flag.StringP("display", "d", "", "The name of file to display")
	var writeFileName = flag.StringP("write", "w", "", "The name of file to write")
	var writeText = flag.StringP("text", "t", "", "The data to write to the file")
	flag.Parse()

	if *makeFileName != "" {
		makeCSV(*makeFileName, *makeForce)
	} 
	if *displayFileName != "" {
		displayFile(*displayFileName)
	}

	if *writeFileName != "" && *writeText != "" {
		writeLine(*writeFileName, *writeText)
	} else if *writeFileName != "" {
		fmt.Println("Please provide data to write to the file")
	}
}