package latex

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type rendererType int64

const (
	RendererTypeHTML rendererType = iota
	RendererTypeMarkdown
)

type BlockRenderer struct {
	RendererType rendererType
}

func (r *BlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindLatexBlock, r.renderLatexBlock)
}

func (r *BlockRenderer) renderLatexBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	switch r.RendererType {
	case RendererTypeHTML:
		return r.renderLatexBlockToHtml(w, source, node, entering)
	case RendererTypeMarkdown:
		return r.renderLatexBlockToMarkdown(w, source, node, entering)
	default:
		panic("unknown renderer type")
	}
}

func (r *BlockRenderer) renderLatexBlockToHtml(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*LatexBlock)
	if entering {
		w.WriteString(`<div class="math block">`)
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			w.Write(line.Value(source))
		}
	} else {
		w.WriteString("</div>")
	}
	return ast.WalkContinue, nil
}

func (r *BlockRenderer) renderLatexBlockToMarkdown(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*LatexBlock)
	if entering {
		w.WriteString(`$$`)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			w.Write(escapeUnderscore(string(segment.Value(source))))
		}
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			w.Write(escapeUnderscore(string(line.Value(source))))
		}
		return ast.WalkSkipChildren, nil
	} else {
		w.WriteString(`$$`)
	}
	return ast.WalkContinue, nil
}

type InlineRenderer struct {
	RendererType rendererType
}

func (r *InlineRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindLatexInline, r.renderLatexInline)
}

func (r *InlineRenderer) renderLatexInline(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	switch r.RendererType {
	case RendererTypeHTML:
		return r.renderLatexInlineToHtml(w, source, node, entering)
	case RendererTypeMarkdown:
		return r.renderLatexInlineToMarkdown(w, source, node, entering)
	default:
		panic("unknown renderer type")
	}
}

func (r *InlineRenderer) renderLatexInlineToHtml(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// n := node.(*LatexInline)
	if entering {
		w.WriteString(`<span class="math inline">`)
		// Inline Latex childrens are all raw text segments
		// so we can just write them out as we don't need to further preprocess
		//
		// for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		// 	segment := c.(*ast.Text).Segment
		// 	w.Write(segment.Value(source))
		// }
		// return ast.WalkSkipChildren, nil
	} else {
		w.WriteString("</span>")
	}
	return ast.WalkContinue, nil
}

func (r *InlineRenderer) renderLatexInlineToMarkdown(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*LatexInline)
	if entering {
		w.WriteString(`$`)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			w.Write(escapeUnderscore(string(segment.Value(source))))
		}
		return ast.WalkSkipChildren, nil
	} else {
		w.WriteString("$")
	}
	return ast.WalkContinue, nil
}

// escape underscores in latex
func escapeUnderscore(s string) []byte {
	// prevent _ in latex to become italic marker
	s = strings.ReplaceAll(s, `\_`, `_`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	// prevent * in latex to become italic marker
	s = strings.ReplaceAll(s, `*`, `\ast `)
	// $\begin{\align*} ... \end{\align*}$
	s = strings.ReplaceAll(s, `\align\ast `, `\align*`)
	// prevent becoming a checkmark
	s = strings.ReplaceAll(s, `\[X]`, `[X]`)
	s = strings.ReplaceAll(s, `[X]`, `\[X]`)
	s = strings.ReplaceAll(s, `\[x]`, `[x]`)
	s = strings.ReplaceAll(s, `[x]`, `\[x]`)
	// prevent escaping of { and }
	s = strings.ReplaceAll(s, `\\{`, `\{`)
	s = strings.ReplaceAll(s, `\{`, `\\{`)
	s = strings.ReplaceAll(s, `\\}`, `\}`)
	s = strings.ReplaceAll(s, `\}`, `\\}`)
	// prevent escaping of spaces
	s = strings.ReplaceAll(s, `\\,`, `\,`)
	s = strings.ReplaceAll(s, `\,`, `\\,`)
	s = strings.ReplaceAll(s, `\\:`, `\:`)
	s = strings.ReplaceAll(s, `\:`, `\\:`)
	s = strings.ReplaceAll(s, `\\;`, `\;`)
	s = strings.ReplaceAll(s, `\;`, `\\;`)
	s = strings.ReplaceAll(s, `\\>`, `\>`)
	s = strings.ReplaceAll(s, `\>`, `\\>`)
	return []byte(s)
}
