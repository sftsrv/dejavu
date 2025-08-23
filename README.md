---
patterns:
  - "^#.*"
---

# Déjà vu

Developers don't read documentation. Just bring it to them instead.

`dejavu` surfaces documentation that may be relevant to your developers alongside the output of your normal CLI commands

## Installation

The project is still in the early stages of development, at the moment it's only possible to install using `go install`

### Using `go install`

```sh
go install github.com/sftsrv/dejavu@latest
```

### Usage

The application can be used by piping any other command on which you'd like to surface documentation, for example `cat myfile.txt`:

```sh
cat myfile.txt | dejavu --docs ./my-docs-dir
```

> `--docs` is a path to the documentation to reference. If omitted this will search all `.md` files in the directory

Each doc should be a markdown file that contains frontmatter with a `patterns` property which is a list of regular expressions on which the document should be shown:

> Patterns are matched on a single line

```md
---
patterns:
  - "^MY_ERROR"
  - "BADBAD"
---

# MY_ERROR or BADBAD Error Handling Doc

These are common errors that we can see when running some tests,
they are triggered if we see output with a line starting with
`MY_ERROR` or any text containing `BADBAD`
```

## Development

Basically, find something to pipe, this can be anything that sends some output, e.g `cat myfile.txt` and pipe that into `go run .`:

```sh
cat myfile.txt | go run .
```

## References

- [Examples for using io.Pipe in go](https://www.zupzup.org/io-pipe-go/index.html)
- [Mastering io.Pipe in go](https://medium.com/@0xgotznit/mastering-io-pipe-in-go-ca8686150b5e)

## Project Roadmap

Some things that I still want to do before considering this project complete:

- [ ] Allow filtering based on type of doc, e.g. errors only
- [ ] Allow docs to specify a summary? User can view the full doc when desired
- [ ] Make it possible to ignore some paths when searching for docs (e.g. `node_modules`)
- [ ] Should support stderr + stdout
