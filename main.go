package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getMessage(line string) string {
	if strings.Contains(line, "func") {
		return "LOOK IT'S A FUNCTION"
	}

	return ""
}

type Doc struct {
	path      string
	contents  string
	familiars []string
}

func parseFile(path string) (Doc, error) {
	return Doc{}, fmt.Errorf("Not implemented")
}

func loadDocs(root string) []Doc {
	var docs []Doc

	filepath.WalkDir(root,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, ".md") {
				doc, docErr := parseFile(s)
				if docErr == nil {
					docs = append(docs, doc)
				}
			}

			return nil
		},
	)

	return docs
}

func main() {
	dir := flag.String("docs", "./", "path to directory containing docs")

	docs := loadDocs(*dir)

	log.Print(docs)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()

		message := getMessage(text)
		if message != "" {
			fmt.Fprintf(os.Stdout, "\r\n%s", message)
		}

		os.Stdout.WriteString("\r\n")
		os.Stdout.Write(scanner.Bytes())
	}
}
