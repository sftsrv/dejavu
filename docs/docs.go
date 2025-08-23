package docs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/goccy/go-yaml"
)

type Doc struct {
	Path     string
	Summary  string
	Rendered string
	Tags     []string
	Patterns []*regexp.Regexp
}

type Frontmatter struct {
	Summary  string   `json:"summary"`
	Tags     []string `json:"tags"`
	Patterns []string `json:"patterns"`
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

func parseFrontmatter(contents string) Frontmatter {
	result := Frontmatter{}

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

func parseFile(path string, summary bool) (Doc, error) {
	file, err := os.ReadFile(path)
	contents := string(file)

	if err != nil {
		return Doc{}, fmt.Errorf("Cannot find the specified file %s", path)
	}

	frontmatter := parseFrontmatter(contents)
	patterns := compilePatterns(frontmatter.Patterns)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	message := clearFrontmatter(contents)
	if summary && frontmatter.Summary != "" {
		message = frontmatter.Summary
	}

	rendered, err := renderer.Render(message)

	if err != nil {
		rendered = contents
	}

	rendered = strings.TrimSpace(rendered)

	return Doc{
		path,
		frontmatter.Summary,
		rendered,
		frontmatter.Tags,
		patterns,
	}, nil

}

func loadDocs(root string, summary bool) []Doc {
	var docs []Doc

	filepath.WalkDir(root,
		func(s string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(s, ".md") {
				doc, docErr := parseFile(s, summary)
				if docErr == nil {
					docs = append(docs, doc)
				}
			}

			return nil
		},
	)

	return docs
}

func intersects(as, bs []string) bool {
	for _, a := range as {
		if slices.Contains(bs, a) {
			return true
		}
	}

	return false
}

func filterDocs(docs []Doc, filter []string) []Doc {
	if len(filter) == 0 {
		return docs
	}

	var filtered []Doc

	for _, doc := range docs {
		if intersects(filter, doc.Tags) {
			filtered = append(filtered, doc)
		}
	}

	return filtered
}

func Load(root string, filter []string, summary bool) []Doc {
	docs := loadDocs(root, summary)

	return filterDocs(docs, filter)
}
