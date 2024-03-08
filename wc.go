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

/*
// function to calculate total number of bytes in a file
func getBytesInFile(fileName string) (int, error) {
	var totalWords = 0
	fileContent, err := getFileContent(fileName)
	if err!=nil {
		return totalWords, err
	}
	totalWords = len(fileContent)
	return totalWords, nil
}
*/

// function to count total number of characters in a file
func getTotalCharacters(fileContent *string) (int) {
	var totalCharacters = 0
	totalCharacters = len(string(*fileContent))
	return totalCharacters 
}

// print results

func main(){
	/*
	// technique to get the present working directory
	pwd, dirErr := os.Getwd()
	if dirErr != nil {
		fmt.Println("Directory parsing error " + dirErr.Error())
	}
	*/
	fmt.Println("Welcome to Word count program!!") 
	var fileToReadForLines = flag.String("l", "", "No of lines in the file, (if present)")
	var fileToReadForWords = flag.String("w", "", "No of words in the file, (if present)")
//	var fileToReadForBytes = flag.String("c", "", "No of bytes in the file, (if present)")
	var fileToReadForCharacters = flag.String("c", "", "No of characters in the file, (if present)")

	flag.Parse()

	// fileContentInBytes : a pointer to file content in bytes, *[]uint8
	// fileContentInBytes, err := getFileContent()
	if *fileToReadForLines != "" {
		fileContentInBytes, err := getFileContent(*fileToReadForLines)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fileContent := string(*fileContentInBytes)
		lines := getLinesInFile(&fileContent)
		fmt.Printf("%8d %s\n", len(*lines), *fileToReadForLines)
	} 
	if *fileToReadForWords != "" {
		fileContentInBytes, err := getFileContent(*fileToReadForWords)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fileContent := string(*fileContentInBytes)
		lines := getLinesInFile(&fileContent)
		words := getWordsFromLines(lines)
		fmt.Printf("%8d %s\n", len(*words), *fileToReadForWords)
	} 
	/*
	if *fileToReadForBytes != "" {
		totalBytes, err := getBytesInFile(*fileToReadForBytes)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Printf("\t%d %v\n", totalBytes, *fileToReadForBytes)
	} 
	*/

	if *fileToReadForCharacters != "" {
		fileContentInBytes, err := getFileContent(*fileToReadForCharacters)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fileContent := string(*fileContentInBytes)
		totalChars := getTotalCharacters(&fileContent)
		fmt.Printf("%8d %s\n", totalChars, *fileToReadForCharacters)
	} 

	// if none of the flags are used then print all three values total number of lines, words and bytes
	if *fileToReadForWords == "" && *fileToReadForLines == "" && *fileToReadForCharacters == "" {
		// if there is atleast one file passed as an argument
		if len(os.Args) > 1 {
			filesToRead := os.Args[1:]
			totalLines, totalWords, totalChars := 0, 0, 0
			for _, file := range filesToRead {
				fileContentInBytes, err := getFileContent(file)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}
				fileContent := string(*fileContentInBytes)
				chars := getTotalCharacters(&fileContent)
				lines := getLinesInFile(&fileContent)
				words := getWordsFromLines(lines)
				fmt.Printf("%8d %8d %8d %s\n", len(*lines), len(*words), chars, file)
				
				// get the total number of lines, words and characters in all the files, if there are atleast 2 files
				if len(filesToRead) > 1 {
					totalLines += len(*lines)
					totalWords += len(*words)
					totalChars += chars
				}
					
			}
			if len(filesToRead) > 1 {
				fmt.Printf("%8d %8d %8d %s", totalLines, totalWords, totalChars, "total") 
			}
		}
	}   
}

/*
// Testing Commands:

./wc -l shakespeare-db/testTextNoRead.txt
./wc -l shakespeare-db/testTextNoRead.txtadsfa
./wc -l shakespeare-db/testText.txt
./wc -l shakespeare-db/
./wc -w shakespeare-db/testTextNoRead.txt
./wc -w shakespeare-db/testTextNoRead.txtadsfa
./wc -w shakespeare-db/testText.txt
./wc -w shakespeare-db/
./wc -c shakespeare-db/testTextNoRead.txt
./wc -c shakespeare-db/testTextNoRead.txtadsfa
./wc -c shakespeare-db/testText.txt
./wc -c shakespeare-db/
*/
