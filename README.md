[![GoDoc](https://godoc.org/jhhgo.us/markdown?status.svg)](https://godoc.org/jhhgo.us/markdown)
# package markdown

`import "jhhgo.us/markdown"`

Package markdown translates plain text with Markdown formatting and YAML front-matter into HTML documents with metadata. Markdown processing is done by package [blackfriday](https://github.com/russross/blackfriday).

Document format is:

    ---
    title: Complete
    date: 2016-05-27T23:45:46.000Z # RFC3339 format
    categories:
      - test
    tags:
      - foo
      - bar
    ---
    # Heading
    The remainder of the document in *Markdown* format.
