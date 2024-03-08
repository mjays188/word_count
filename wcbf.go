package main

import(
	"fmt"
	"flag"
	// file reading packages
	"os"
	"io"
//	"reflect"
	"strings"
	"errors"
)

// function to get the content of the file, or through relevant error
func getFileContent(fileName string) (*[]uint8, error) {
	/*
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("getFileContent %v\n", reflect.TypeOf(err))	// *fs.PathError
		fmt.Println(err)
	}
	return string(content), nil
	*/
	var content []uint8
	// first of all check if the a file/directory with the given fileName exists
	file, err := os.Open(fileName)
	defer file.Close() 
	if err != nil {
		if os.IsNotExist(err) {
			//fmt.Println("IsNotExist " + err.Error())
			return &content, errors.New(fmt.Sprintf("./wc: %s: open: No such file or directory\n", fileName))
		} else if os.IsPermission(err) {
			//fmt.Println("IsPermission " + err.Error())
			return &content, errors.New(fmt.Sprintf("./wc: %s: open: Permission denied\n", fileName))
		} else {
			//fmt.Println("Error opening file " + err.Error())
			return &content, errors.New(fmt.Sprintf("./wc: %s: %v\n", fileName, err.Error()))
		}
		return &content, err
	}

	// check if returned file is a directory or not
	fileInfo, err := file.Stat()
	if err!=nil {
		//fmt.Println("Error fetching fileInfo: " + err.Error())		
		return &content, errors.New(fmt.Sprintf("./wc: %s: %v\n", fileName, err.Error()))
	}
	if fileInfo.Mode().IsDir() {
		//fmt.Println(fileName + " is a directory ")
		return &content, errors.New(fmt.Sprintf("./wc: %s: read: Is a directory\n", fileName))
	}

	// read the contents of the file
	content, err = io.ReadAll(file)
	if err!=nil {
		return &content, errors.New(fmt.Sprintf("./wc: %s: %v\n", fileName, err.Error()))
	}
	//fmt.Println(content)

	return &content, nil
}


func getLinesInFile(fileContent *string) (*[]string) {
	var lines []string
	// split the string with '\n' as delimiter
	lines = strings.Split(*fileContent, "\n")

	// to avoid miscalculating the last empty line
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return &lines
}


// function to get the total words in the file
func getWordsFromLines(lines *[]string) (*[]string) {
	var words []string
	
	// split each line with whitespace as delimiter and accumulate all the words
	for _, line := range *lines {
		wordsInLine := strings.Split(line, " ")
		words = append(words, wordsInLine...) 
	}

	return &words	
}

// function to count total number of characters in a file
func getCharacterCount(fileContent *string) (int) {
	var totalCharacters = 0
	totalCharacters = len(string(*fileContent))
	return totalCharacters 
}

// print results
func printResult(l, w, c bool, caption string, lines, words, chars int) {
	if l {
		fmt.Printf("%8d", lines)
	}
	if w {
		if l {
			fmt.Printf(" ")
		}
		fmt.Printf("%8d", words)
	}
	if c {
		if l || w {
			fmt.Printf(" ")
		}
		fmt.Printf("%8d", chars)
	}
	fmt.Printf(" %s\n", caption)
}

func main(){
	/*
	// technique to get the present working directory
	pwd, dirErr := os.Getwd()
	if dirErr != nil {
		fmt.Println("Directory parsing error " + dirErr.Error())
	}
	*/

	// check if multiple flags are passed as a combination, for ex:
	// ./wc -l -w <file1> <file2> <file3> ...
	// ./wc -c -w <file1> <file2> <file3> ...
	// ./wc -l -c <file1> <file2> <file3> ...
	// ./wc -l -c -w <file1> <file2> <file3> ...
	var getTotalLines = flag.Bool("l", false, "No of lines in the file, (if present)")
	var getTotalWords = flag.Bool("w", false, "No of words in the file, (if present)")
	var getTotalCharacters = flag.Bool("c", false, "No of characters in the file, (if present)")

	flag.Parse()

	// if none of the flags are used then print all three values total number of lines, words and bytes
	if *getTotalLines || *getTotalWords || *getTotalCharacters {
		// if there is atleast one file passed as an argument
		if len(flag.Args()) > 0 {
			filesToRead := flag.Args()
			totalLines, totalWords, totalChars := 0, 0, 0
			for _, file := range filesToRead {
				fileContentInBytes, err := getFileContent(file)
				if err != nil {
					fmt.Printf("%s", err.Error())
					if len(filesToRead) == 1 {
						os.Exit(0)
					}
					continue
				}
				fileContent := string(*fileContentInBytes)
				chars := getCharacterCount(&fileContent)
				lines := getLinesInFile(&fileContent)
				words := getWordsFromLines(lines)
				printResult(*getTotalLines, *getTotalWords, *getTotalCharacters, file, len(*lines), len(*words), chars)
				
				// get the total number of lines, words and characters in all the files, if there are atleast 2 files
				if len(filesToRead) > 1 {
					totalLines += len(*lines)
					totalWords += len(*words)
					totalChars += chars
				}
			}
			if len(filesToRead) > 1 {
				printResult(*getTotalLines, *getTotalWords, *getTotalCharacters, "Total", totalLines, totalWords, totalChars) 
			}
		} else {
			// implement the code to take input from stdin and get all 3 details for entered text on pressing ctrl+d
		}
	}   
}
