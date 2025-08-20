# Déjà vu

Developers don't read documentation. Just bring it to them instead.

`dejavu` surfaces documentation that may be relevant to your developers within outputs of CLI commands

## Usage

## Development

Basically, find something to pipe, this can be anything that sends some output, e.g `cat myfile.txt` and pipe that into `go run .`:

```sh
cat myfile.txt | go run .
```

## References

- [Examples for using io.Pipe in go](https://www.zupzup.org/io-pipe-go/index.html)
- [Mastering io.Pipe in go](https://medium.com/@0xgotznit/mastering-io-pipe-in-go-ca8686150b5e)
