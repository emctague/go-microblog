package main

import (
	"bufio"
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
)

func SeparateFrontMatter(in io.Reader, structure interface{}) ([]byte, error) {
	frontMatter := new(bytes.Buffer)
	frontMatterOut := bufio.NewWriter(frontMatter)

	remaining := new(bytes.Buffer)
	remainingOut := bufio.NewWriter(remaining)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "---" {
			break
		}
		_, _ = frontMatterOut.WriteString(text + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := frontMatterOut.Flush(); err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(frontMatter.Bytes(), structure); err != nil {
		return nil, err
	}

	for scanner.Scan() {
		_, _ = remainingOut.Write(scanner.Bytes())
		_ = remainingOut.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := remainingOut.Flush(); err != nil {
		return nil, err
	}

	return remaining.Bytes(), nil
}
