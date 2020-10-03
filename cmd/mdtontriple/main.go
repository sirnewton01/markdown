package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// This prints AST of parsed markdown document.
// Usage: printast <markdown-file>

func usageAndExit() {
	fmt.Printf("Usage: printast <markdown-file>\n")
	os.Exit(1)
}

func main() {
	nFiles := len(os.Args) - 1
	if nFiles < 1 {
		usageAndExit()
	}
	for i := 0; i < nFiles; i++ {
		fileName := os.Args[i+1]
		d, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open '%s', error: '%s'\n", fileName, err)
			continue
		}
		exts := parser.NoIntraEmphasis |
		parser.Tables |
		parser.FencedCode |
		parser.Autolink |
		parser.Strikethrough |
		parser.SpaceHeadings |
		parser.Attributes

		p := parser.NewWithExtensions(exts)
		doc := markdown.Parse(d, p)

		bnctr := 1
		hs := ""

		ast.WalkFunc(doc, func(n ast.Node, entering bool) ast.WalkStatus {
			if h, ok := n.(*ast.Heading); ok && !entering && len(h.Children) == 1 {
				t := h.Children[0].AsLeaf()
				l := fmt.Sprintf("\"%s\"", string(t.Literal))

				lp := "<http://www.w3.org/2000/01/rdf-schema#label>"

				if h.Level == 1 {
					hs = "<>"
				} else if h.Attribute != nil {
					hs = fmt.Sprintf("<#%s>", h.Attribute.ID)
				} else {
                                        hs = fmt.Sprintf("_:bn%d", bnctr)
					bnctr++
				}

				fmt.Printf("%s %s %s.\n", hs, lp, l)

			}

			if li, ok := n.(*ast.ListItem); ok && !entering && len(li.Children) == 1 {
				p := li.Children[0].AsContainer()
				if l, ok := p.Children[1].(*ast.Link); ok {
					r := string(l.Destination)

					t := p.Children[2].AsLeaf()
					if len(string(t.Literal)) == 0 {
						fmt.Printf("%s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <%s>.\n", hs, r)
					} else {
						fmt.Printf("%s <%s> \"%s\".\n", hs, r, strings.TrimSpace(string(t.Literal)))
					}
				}
			}
			return ast.GoToNext
		})
	}
}
