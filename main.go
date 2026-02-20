package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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

func runDejavu(docs *[]docs.Doc, streams DejavuStreams) {
	scanner := bufio.NewScanner(streams.in)
	for scanner.Scan() {
		bytes := scanner.Bytes()

		message := getMessage(docs, bytes)
		if message != "" {
			fmt.Fprintf(streams.out, "\r\n%s", message)
		}

		streams.out.WriteString("\r\n")
		streams.out.Write(bytes)
	}
}

type DejavuStreams struct {
	in  io.Reader
	out *os.File
}

func createStdinStream(stdin, stdout *os.File) (DejavuStreams, error) {
	return DejavuStreams{stdin, stdout}, nil
}

func createCommandStream(command string, args []string, in, out *os.File) (DejavuStreams, error) {
	cmd := exec.Command(command, args...)

	cmd.Stdin = in
	// Assigning cmd.Stdin = os.Stdout would enable color output here
	// which can probably be achieved using some type of io.MultiWriter
	// but interactive commands with merged output is already weird enough
	// so I don't know if it's even worth doing

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return DejavuStreams{}, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return DejavuStreams{}, err
	}

	reader := io.MultiReader(stdout, stderr)

	go cmd.Start()

	return DejavuStreams{reader, out}, nil

}

const usage = `dejavu

dejavu surfaces documentation that may be relevant to your developers alongside the output of your normal CLI commands

dejavu is used a pipe as follows:

$ my special command | dejavu <...flags>

or for interactive commands, you can use the following

$ dejavu -c "my special command"

the following flags may be provided when running dejavu. provided flags will override the matching value in the config file:

`

func main() {
	commandFlag := flag.String("c", "", "command to run")
	pathFlag := flag.String("path", "./dejavu.config.json", "path to config file")
	docsFlag := flag.String("docs", "", "path to directory containing docs")
	tagsFlag := flag.String("tags", "", "filter docs by those including the given tags")
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
		Tags:    *tagsFlag,
		Summary: *summaryFlag,
	}

	config := config.Load(flags)

	docs := docs.Load(config.Docs, config.Tags, config.Summary)

	var streams DejavuStreams
	var err error

	if len(*commandFlag) > 1 {
		parts := strings.Fields(*commandFlag)
		name := parts[0]
		args := parts[1:]
		streams, err = createCommandStream(name, args, os.Stdin, os.Stdout)
	} else {
		streams, err = createStdinStream(os.Stdin, os.Stdout)
	}

	if err != nil {
		panic(err)
	}

	runDejavu(&docs, streams)
}
