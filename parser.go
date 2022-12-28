package latex

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Parser struct{}

var _ parser.InlineParser = &Parser{}

func (p *Parser) Trigger() []byte {
	return []byte{'$'}
}

func (p *Parser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, seg := block.PeekLine()

	if !bytes.HasPrefix(line, []byte{'$'}) {
		return nil
	}

	if len(line) < 2 || line[1] == ' ' {
		// $ asdf -> not math
		return nil
	}

	if line[1] != '$' {
		// inline math
		stop := bytes.Index(line[1:], []byte{'$'})
		if stop < 0 {
			return nil // must close on the same line
		}

		n := &LatexInline{}
		seg = text.NewSegment(seg.Start+1, seg.Start+stop+1)
		n.AppendChild(n, ast.NewRawTextSegment(seg))
		block.Advance(stop + 2)
		return n
	}

	// $$ -> block math
	stop := bytes.Index(line[2:], []byte{'$', '$'})
	if stop < 0 {
		return nil
	}

	n := &LatexBlock{}
	seg = text.NewSegment(seg.Start+2, seg.Start+stop+2)
	n.AppendChild(n, ast.NewRawTextSegment(seg))
	block.Advance(stop + 4)
	return n
}
