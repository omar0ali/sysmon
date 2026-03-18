package internal

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const PROC_DIR = "/proc"

func Must[T any](i T, e error) T {
	if e != nil {
		panic(e)
	}
	return i
}

func open(path string) (*os.File, *bufio.Scanner, error) {
	if runtime.GOOS != "linux" {
		return nil, nil, fmt.Errorf("this tool only works on Linux")
	}

	if _, err := os.Stat("/proc"); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("/proc is not available (procfs not mounted)")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(file)
	return file, scanner, nil
}

// this open will close the file automatically
func OpenWithScanner(path string, f func(scanner *bufio.Scanner)) error {
	file, scanner, err := open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	f(scanner)
	if err := scanner.Err(); err != nil {
		return err
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
