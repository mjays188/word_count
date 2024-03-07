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
func getFileContent(fileName string) ([]uint8, error) {
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
			return content, errors.New(fmt.Sprintf("./wc: %s: open: No such file or directory\n", fileName))
		} else if os.IsPermission(err) {
			//fmt.Println("IsPermission " + err.Error())
			return content, errors.New(fmt.Sprintf("./wc: %s: open: Permission denied\n", fileName))
		} else {
			//fmt.Println("Error opening file " + err.Error())
			return content, errors.New(fmt.Sprintf("./wc: %s: %v\n", fileName, err.Error()))
		}
		return content, err
	}

	// check if returned file is a directory or not
	fileInfo, err := file.Stat()
	if err!=nil {
		//fmt.Println("Error fetching fileInfo: " + err.Error())		
		return content, errors.New(fmt.Sprintf("./wc: %s: %v\n", fileName, err.Error()))
	}
	if fileInfo.Mode().IsDir() {
		//fmt.Println(fileName + " is a directory ")
		return content, errors.New(fmt.Sprintf("./wc: %s: read: Is a directory\n", fileName))
	}

	// read the contents of the file
	content, err = io.ReadAll(file)
	if err!=nil {
		return content, errors.New(fmt.Sprintf("./wc: %s: %v\n", fileName, err.Error()))
	}
	//fmt.Println(content)

	return content, nil
}


func getLinesInFile(fileName string) (*[]string, error) {
	var lines []string
	// check whether the give fileName is a directory or a file 
	content, err := getFileContent(fileName)
	if err != nil {
		return &lines, err
	}
	// split the string with '\n' as delimiter
	lines = strings.Split(string(content), "\n")
	
	/*
	for _, line := range lines {
		fmt.Println(line)
	}
	*/

	// to avoid miscalculating the last empty line
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return &lines, nil
}

// function to get the total words in the file
func getWordsInFile(fileName string) (*[]string, error) {
	var words []string
	
	lines, err := getLinesInFile(fileName)
	if err != nil {
		//fmt.Println(err.Error())
		return &words, err	
	}
	//fmt.Printf("\t%d %v\n", len(*lines), *fileName)

	// split each line with whitespace as delimiter and accumulate all the words
	for _, line := range *lines {
		wordsInLine := strings.Split(line, " ")
		words = append(words, wordsInLine...) 
	}
	return &words, nil	
}

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

// function to count total number of characters in a file
func getTotalCharacters(fileName string) (int, error) {
	var totalCharacters = 0
	fileContent, err := getFileContent(fileName)
	if err != nil {
		return totalCharacters, err
	}
	totalCharacters = len(string(fileContent))
	return totalCharacters, nil
}

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
	var fileToReadForBytes = flag.String("c", "", "No of bytes in the file, (if present)")
	var fileToReadForCharacters = flag.String("m", "", "No of characters in the file, (if present)")

	flag.Parse()

	if *fileToReadForLines != "" {
		lines, err := getLinesInFile(*fileToReadForLines)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Printf("\t%d %v\n", len(*lines), *fileToReadForLines)
	} 
	if *fileToReadForWords != "" {
		words, err := getWordsInFile(*fileToReadForWords)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Printf("\t%d %v\n", len(*words), *fileToReadForWords)
	} 
	if *fileToReadForBytes != "" {
		totalBytes, err := getBytesInFile(*fileToReadForBytes)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Printf("\t%d %v\n", totalBytes, *fileToReadForBytes)
	} 

	if *fileToReadForCharacters != "" {
		totalChars, err := getTotalCharacters(*fileToReadForCharacters)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Printf("\t%d %v\n", totalChars, *fileToReadForCharacters)
	} 
}
