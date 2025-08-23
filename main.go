package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/sftsrv/dejavu/config"
	"github.com/sftsrv/dejavu/docs"
)

func renderDoc(doc docs.Doc) string {
	return doc.Path + "\n" + doc.Rendered
}

var title = lipgloss.NewStyle().Render("Déjà vu")

func getMessage(docs *[]docs.Doc, line []byte) string {
	message := ""

	for _, doc := range *docs {
		for _, pattern := range doc.Patterns {
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

func processStream(in, out *os.File, docs *[]docs.Doc) {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		bytes := scanner.Bytes()

		message := getMessage(docs, bytes)
		if message != "" {
			fmt.Fprintf(out, "\r\n%s", message)
		}

		out.WriteString("\r\n")
		out.Write(bytes)
	}

}

const usage = `dejavu

dejavu surfaces documentation that may be relevant to your developers alongside the output of your normal CLI commands

dejavu is used a pipe as follows:

$ my special command | dejavu <...flags>

the following flags may be provided when running dejavu. provided flags will override the matching value in the config file:

`

func main() {
	pathFlag := flag.String("path", "./dejavu.config.json", "path to config file")
	docsFlag := flag.String("docs", "", "path to directory containing docs")
	typesFlag := flag.String("types", "", "limit the types of docs to include")
	summaryFlag := flag.Bool("summary", false, "only show summary of doc")
	helpFlag := flag.Bool("help", false, "show help menu")

	flag.Parse()

	if *helpFlag {
		fmt.Print(usage)
		flag.Usage()
		return
	}

	flags := config.Flags{
		Path:    *pathFlag,
		Docs:    *docsFlag,
		Types:   *typesFlag,
		Summary: *summaryFlag,
	}

	config := config.Load(flags)

	docs := docs.Load(config.Docs, config.Types, config.Summary)

	processStream(os.Stdin, os.Stdout, &docs)
}
