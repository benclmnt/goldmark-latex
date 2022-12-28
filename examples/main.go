package main

import (
	"bytes"

	latex "github.com/benclmnt/goldmark-latex"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func main() {
	buf := &bytes.Buffer{}
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithInlineParsers(util.Prioritized(&latex.Parser{}, 499)),
		),
		goldmark.WithRendererOptions(
			renderer.WithNodeRenderers(util.Prioritized(&latex.BlockRenderer{
				RendererType: latex.RendererTypeMarkdown,
			}, 502)),
			renderer.WithNodeRenderers(util.Prioritized(&latex.InlineRenderer{
				RendererType: latex.RendererTypeMarkdown,
			}, 501)),
		),
	)
	md.Convert([]byte(`Example: $\mathcal{M}_{n \times n}(\mathbf{F}) \cong_\mathbb{F} \mathbb{F}^{mn}$, $$\mathcal{P}_n(\mathbb{F}) \cong_\mathbb{F} \mathbb{F}^{n + 1}$$, $\mathbb{C}^n \cong_\mathbb{R} \mathbb{R}^{2n}$`), buf)
}
