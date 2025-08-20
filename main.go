package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	yaml "github.com/goccy/go-yaml"
)

func renderDoc(doc doc) string {
	return doc.path + "\n" + doc.rendered
}

var title = lipgloss.NewStyle().Render("Déjà vu")

func getMessage(docs *[]doc, line []byte) string {
	message := ""

	for _, doc := range *docs {
		for _, pattern := range doc.patterns {
			if pattern.Match(line) {
				message += renderDoc(doc)
			}
		}
	}

	if message == "" {
		return ""
	}

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("72")).
		Border(lipgloss.NormalBorder()).
		Render("  >", title, " ", message)

}

type doc struct {
	path     string
	rendered string
	patterns []*regexp.Regexp
}

type frontmatter struct {
	Patterns []string `yaml:"patterns"`
}

func compilePatterns(patterns []string) []*regexp.Regexp {
	result := []*regexp.Regexp{}

	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			result = append(result, re)
		}
	}

	return result
}

func parseFrontmatter(contents string) frontmatter {
	result := frontmatter{}

	after, start := strings.CutPrefix(contents, "---")

	if !start {
		return result
	}

	lines := strings.Lines(after)
	frontmatter := ""

	for line := range lines {
		if strings.HasPrefix(line, "---") {
			break
		}

		frontmatter += line
	}

	yaml.Unmarshal([]byte(frontmatter), &result)

	return result
}

func clearFrontmatter(contents string) string {
	re := regexp.MustCompile("(?m)^---(.|\n|\r)*?---")
	return strings.TrimSpace(re.ReplaceAllString(contents, ""))
}

func parseFile(path string) (doc, error) {
	file, readErr := os.ReadFile(path)
	contents := string(file)

	if readErr != nil {
		return doc{}, fmt.Errorf("Cannot find the specified file %s", path)
	}

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	rendered, renderErr := renderer.Render(clearFrontmatter(contents))

	if renderErr != nil {
		rendered = contents
	}

	frontmatter := parseFrontmatter(contents)
	patterns := compilePatterns(frontmatter.Patterns)

	return doc{
		path,
		rendered,
		patterns,
	}, nil

}

func loadDocs(root string) []doc {
	var docs []doc

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
	flag.Parse()

	docs := loadDocs(*dir)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		bytes := scanner.Bytes()

		message := getMessage(&docs, bytes)
		if message != "" {
			fmt.Fprintf(os.Stdout, "\r\n%s", message)
		}

		os.Stdout.WriteString("\r\n")
		os.Stdout.Write(bytes)
	}
}
