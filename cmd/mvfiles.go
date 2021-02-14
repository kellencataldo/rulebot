package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var rulebooks = map[string]bool{
	"core":  true,
	"apg":   true,
	"aoepg": true,
}

func validatePrefix(prefix string) bool {
	_, ok := rulebooks[prefix]
	return ok
}

func validateDir(rbDir string) bool {

	if stat, err := os.Stat(rbDir); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func determineSP(f os.FileInfo) (int, bool) {

	if f.IsDir() {
		fmt.Println("Unable to determine start page from nested directory")
		return -1, false
	}

	var startPage int
	var endPage int
	var postfix string
	_, err := fmt.Sscanf(f.Name(), "PZO2101 %d-%d %s", &startPage, &endPage, &postfix)
	if err != nil {
		fmt.Printf("ugh Error scanning file name: %s, error: %s\n", f.Name(), err)
		return -1, false
	}

	return startPage, true
}

func iterateFiles(startPage int, prefix, extension string) {

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Printf("Error opening directory %s\n", err)
		return
	} else if len(files) == 0 {
		fmt.Printf("Empty directory specified\n")
		return
	}

	if startPage == -1 {
		fmt.Println("Determining start page from file name")
		calculatedSP, ok := determineSP(files[0])
		if !ok {
			fmt.Println("Unable to determine start page")
			return
		}

		startPage = calculatedSP
	}

	for _, f := range files {

		if f.IsDir() {
			fmt.Printf("Error: entry %s is a directory, skipping\n", f.Name())
			continue
		}

		newName := prefix + strconv.Itoa(startPage) + extension
		fmt.Printf("Moving %s -> %s\n", f.Name(), newName)
		err := os.Rename(f.Name(), newName)
		if err != nil {
			fmt.Printf("Error renaming file: %s, error: %s\n", f.Name(), err)
		}

		startPage += 1
	}
}

func main() {

	rbPrefix := flag.String("rulebook", "core", "set rulebook prefix")
	rbDir := flag.String("dir", ".", "set rulebook directory")
	rbStartPageFlag := flag.Int("start", -1, "set rulebok start page")
	rbExtension := flag.String("ext", ".png", "set file extension")
	flag.Parse()

	if !validatePrefix(*rbPrefix) {
		fmt.Printf("%s is not a valid rulebook prefix\n", *rbPrefix)
		return
	}

	fmt.Printf("File prefix will be: %s\n", *rbPrefix)

	if !validateDir(*rbDir) {
		fmt.Printf("Unable to find directory path, or path does not represent directory: %s\n", *rbDir)
		return
	}

	if err := os.Chdir(*rbDir); err != nil {
		fmt.Printf("Unable to change to directory: %s, error: %s\n", *rbDir, err)
		return
	}

	iterateFiles(*rbStartPageFlag, *rbPrefix, *rbExtension)
}
