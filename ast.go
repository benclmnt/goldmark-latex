package latex

import "github.com/yuin/goldmark/ast"

type LatexInline struct {
	ast.BaseInline
}

var _ ast.Node = &LatexInline{}

var KindLatexInline = ast.NewNodeKind("LatexInline")

func (n *LatexInline) Kind() ast.NodeKind {
	return KindLatexInline
}

func (n *LatexInline) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, nil, nil)
}

type LatexBlock struct {
	ast.BaseBlock
}

var _ ast.Node = &LatexBlock{}

var KindLatexBlock = ast.NewNodeKind("LatexBlock")

func (n *LatexBlock) Kind() ast.NodeKind {
	return KindLatexBlock
}

func (n *LatexBlock) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, nil, nil)
}
