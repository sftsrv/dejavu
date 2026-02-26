# Déjà vu

Developers don't search for documentation. Just bring it to them instead.

`dejavu` surfaces documentation that may be relevant to your developers alongside the output of your normal CLI commands

![dejavu in action](./images/demo.gif)

## Installation

The project is still in the early stages of development, at the moment it's only possible to install using `go install`

### `go install`

```sh
go install github.com/sftsrv/dejavu@latest
```

## Usage

> [!NOTE]  
> Welcome! If you're using this then consider yourself now a bug hunter. Let me know if you find any issues so that we can make this better for everyone!

### As a Stream

The application can be used by piping any other command on which you'd like to surface documentation, for example `cat myfile.txt`:

```sh
cat myfile.txt | dejavu
```

> [!TIP]
> Since `dejavu` works with the Stdin stream only, in some cases you may need to merge Stin and Stderr into a single stream so that it can be used as expected, doing will depend on the shell you're using

### As a Command

`dejavu` can also execute the command directly while connecting Stdin to the command's Stdin. This can be done by passing the entire command to run to `dejavu`

```sh
dejavu -c "some interactive command"
```

### As a Query

Or if you just have a message as a string that you'd like to get some matches on - you can run it with:

```sh
dejavu -q "my error message"
```

> This is equivalent to something like `echo "my error message" | dejavu` but is a bit less mysterious

## Documentation Structure

Each doc should be a markdown file that contains frontmatter with a `patterns` property which is a list of regular expressions on which the document should be shown:

> Patterns are matched on a single line

```md
---
# list of patterns to match for this doc to be shown
patterns:
  - "^MY_ERROR"
  - "BADBAD"
# list of tags that apply to this doc
tags:
  - error
  - setup
# summary to be shown if dejavu is running in `summary` mode
summary: This is a bad error
---

# MY_ERROR or BADBAD Error Handling Doc

These are common errors that we can see when running some tests,
they are triggered if we see output with a line starting with
`MY_ERROR` or any text containing `BADBAD`
```

## Configuration

`dejavu` can be configured by using a `dejavu.config.json` file with the following structure:

```json
{
  "docs": "./faq",
  "summary": true,
  "tags": [
    "warning",
    "error"
  ]
}
```

> The config file can also be set with the `--config` flag.

Values from the configuration may also be overidden using the matching command line arguments. These can be found by using `dejavu --help`


## Development

Basically, find something to pipe, this can be anything that sends some output, e.g `cat myfile.txt` and pipe that into `go run .`:

```sh
cat myfile.txt | go run .
```

## References

- [Examples for using io.Pipe in go](https://www.zupzup.org/io-pipe-go/index.html)
- [Mastering io.Pipe in go](https://medium.com/@0xgotznit/mastering-io-pipe-in-go-ca8686150b5e)
