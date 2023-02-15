package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadLines(file string) ([]string, error) {
	return ReadSplitter(file, '\n')
}

func ReadSplitter(file string, splitter byte) (lines []string, err error) {
	fin, err := os.Open(file)
	if err != nil {
		return
	}

	r := bufio.NewReader(fin)
	for {
		line, err := r.ReadString(splitter)
		if err == io.EOF {
			break
		}
		line = strings.Replace(line, string(splitter), "", -1)
		lines = append(lines, line)
	}
	return
}

func ReadFileToString(file string) (content string, err error) {
	fin, err := os.Open(file)
	if err != nil {
		return content, err
	}

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if scannerError := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", scannerError)
		return content, err
	}

	return content, err
}
