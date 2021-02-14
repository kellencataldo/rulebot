package main

import (
	"flag"
	"fmt"
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

func iterateFiles(endPage int, prefix, ext string) {

	for i := 1; i <= endPage; i++ {
		filename := prefix + strconv.Itoa(i) + ext
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("File missing: %s\n", filename)
		}
	}

	fmt.Println("Iteration complete")
}

func main() {

	rbPrefix := flag.String("rulebook", "core", "set rulebook prefix")
	rbDir := flag.String("dir", ".", "set rulebook directory")
	rbExtension := flag.String("ext", ".png", "set rulebook file extension")
	rbEndPageFlag := flag.Int("end", -1, "set rulebok end page")
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

	iterateFiles(*rbEndPageFlag, *rbPrefix, *rbExtension)
}
