---
# list of patterns to match for this doc to be shown
patterns:
  - "error"
# list of tags that apply to this doc
tags:
  - error
  - example
summary: This is an example error
---

# Example Error

This is an example error that just lights up on every use of "error".

The easiest way to trigger this is by running `dejavu` with something that has the word `error` in it, for example:

```sh
echo "Hello this is an error" | go run .
```

Or

```sh
go run . -c "echo 'Hello this is an error'"
```
