package main

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	flag "github.com/spf13/pflag"
)

func makeCSV(fileName string, force bool) {

	/*
		Function to create a CSV

		Parameters:
			fileName (string): The name of the file to create
			force (bool): If true, overwrite the file if it already exists
	*/

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
	writer.Write([]string{"ID", "Task", "Time"})
	writer.Flush()
	defer file.Close()
	fmt.Println("File created successfully")
}

func listFile(fileName string) {
	/*
		Function to list the contents of a CSV file

		Parameters:
			fileName (string): The name of the file to display
	*/
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

	// Writing to file
	fmt.Println("Listing file: ", fileName)
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
	for _, row := range data {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}

func writeLine(fileName string, text string) {
	/*
		Function to write a line to a CSV file

		Parameters:
			fileName (string): The name of the file to write to
			text (string): The text to write to the file
	*/
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	id := make([]byte, 4)
	rand.Read(id)
	idString := hex.EncodeToString(id)
	if err != nil {
		panic(err)
	}
	data := []string{idString, text, time.Now().Format("2006-01-02 15:04:05")}
	writer := csv.NewWriter(file)
	writer.Write(data)
	writer.Flush()
}

func deleteTask(filename string, id string) {
	/*
		Function to delete a task from a CSV file

		Parameters:
			fileName (string): The name of the file to delete from
			id (string): The ID of the task to delete
	*/
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var new_data [][]string
	for _, row := range data {
		if row[0] != id {
			new_data = append(new_data, row)
		}
	}

	// Empty the file and go to start of the file
	file.Truncate(0)
	file.Seek(0, 0)

	// Write the new data to the file
	writer := csv.NewWriter(file)
	writer.WriteAll(new_data)
	writer.Flush()
}

func main() {
	var makeFileName = flag.StringP("make", "m", "", "The name of file to make")
	var makeFileNameForce = flag.StringP("make-force", "M", "", "The name of file to overwrite")
	var listFileName = flag.StringP("list", "l", "", "The name of file to display")
	var writeArgs = flag.StringP("write", "w", "", "The name of file follwed by the text to write")
	var deleteArgs = flag.StringP("delete", "d", "", "The name of file to delete followed by ID")
	flag.Parse()

	// To create a file
	if *makeFileName != "" {
		makeCSV(*makeFileName, false)
	}
	// To create a file and overwrite if it already exists
	if *makeFileNameForce != "" {
		makeCSV(*makeFileNameForce, true)
	}
	// To display the contents of a file
	if *listFileName != "" {
		listFile(*listFileName)
	}

	// To write to a file
	if *writeArgs != "" {
		args := flag.Args()

		if len(args) == 0 {
			fmt.Println("Provide text to write")
			return
		}

		if len(args) > 1 {
			fmt.Println("Received more than one arguement. Use double quotes (\") if your text contains space.")
			return
		}

		var text string = args[0]
		writeLine(*writeArgs, text)
	}

	// To delete a line from a file
	if *deleteArgs != "" {
		args := flag.Args()

		if len(args) == 0 {
			fmt.Println("Provide ID to delete")
			return
		}

		if len(args) > 1 {
			fmt.Println("Supports one ID at a time")
			return
		}

		var id string = args[0]
		deleteTask(*deleteArgs, id)
	}
}
