package helper

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Must[T any](i T, e error) T {
	if e != nil {
		panic(e)
	}
	return i
}

// this open needs to close the file later
func Open(path string) (*os.File, *bufio.Scanner) {
	file := Must(os.Open(path))
	scanner := bufio.NewScanner(file)
	return file, scanner
}

// this open will close the file automatically
func OpenScanner(path string, f func(scanner *bufio.Scanner)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	f(scanner)
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return nil
}

func ParseUint(val string) uint64 {
	if val == "" {
		return 0
	}
	num := strings.Fields(val)[0]
	v, _ := strconv.ParseUint(num, 10, 64)
	return v
}
