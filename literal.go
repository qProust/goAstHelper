package asth

import (
	"fmt"
	"go/ast"
	"go/token"
)

type (
	Literal interface {
		Rvalue
		asthLiteralExpr() ast.Expr
	}

	BasicLiteral struct {
		*BaseRvalue
	}

	Identifier BasicLiteral

	StructLiteral struct {
		expr *ast.CompositeLit
	}
	MapLiteral struct {
		expr *ast.CompositeLit
	}
	SliceLit struct {
		expr *ast.CompositeLit
	}
)

var (
	Nil         = &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.IDENT, Value: "nil"}}}
	True        = &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.IDENT, Value: "true"}}}
	False       = &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.IDENT, Value: "false"}}}
	EmptyString = NewStringLiteral("")
)

func NewIdent(name string) *Identifier {
	return &Identifier{&BaseRvalue{expr: ast.NewIdent(name)}}
}
func NewStringLiteral(v string) *BasicLiteral {
	return &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.STRING, Value: "`" + v + "`"}}}
}
func NewBackquoteStringLiteral(v string) *BasicLiteral {
	return &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.STRING, Value: "`" + v + "`"}}}
}
func NewBoolLiteal(b bool) *BasicLiteral {
	return &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.IDENT, Value: fmt.Sprintf("%v", b)}}}
}
func NewIntLiteral(v int) *BasicLiteral {
	return &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", v)}}}
}
func NewUintLiteral(v uint) *BasicLiteral {
	return &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%d", v)}}}
}
func NewFloatLiteral(v float64) *BasicLiteral {
	return &BasicLiteral{&BaseRvalue{&ast.BasicLit{Kind: token.FLOAT, Value: fmt.Sprintf("%f", v)}}}
}
func (l *BasicLiteral) asthLiteralExpr() ast.Expr { return l.expr }
func (l *Identifier) asthLValue() ast.Expr        { return l.expr }

// SLICE

func NewSliceLit() *SliceLit {
	return NewTypedSliceLit(nil)
}

func NewTypedSliceLit(typ Type) *SliceLit {
	s := &SliceLit{
		&ast.CompositeLit{
		},
	}
	if typ != nil {
		s.expr.Type = typ.asthType()
	}
	return s
}

func (l *SliceLit) Append(ls ...Rvalue) *SliceLit {
	if ls == nil {
		return l
	}
	for _, v := range ls {
		if v == nil {
			continue
		}
		l.expr.Elts = append(l.expr.Elts, v.asthRValue())
	}
	return l
}

func (l *SliceLit) asthLiteralExpr() ast.Expr { return l.expr }
func (l *SliceLit) asthRValue() ast.Expr      { return l.expr }

// MAP
func NewMapLiteral(keyType *TypeRef, valType *TypeRef) *MapLiteral {
	return &MapLiteral{
		&ast.CompositeLit{
			Type: &ast.MapType{Key: keyType.node, Value: valType.node},
		},
	}
}
func NewTypedMapLiteral(customType Type) *MapLiteral {
	return &MapLiteral{
		&ast.CompositeLit{
			Type: customType.asthType(),
		},
	}
}
func (l *MapLiteral) AddEntry(key Literal, val Rvalue) *MapLiteral {
	l.expr.Elts = append(l.expr.Elts, &ast.KeyValueExpr{
		Key:   key.asthLiteralExpr(),
		Value: val.asthRValue(),
	})
	return l
}
func (l *MapLiteral) AddEntries(es map[Literal]Literal) *MapLiteral {
	if es == nil {
		return l
	}
	for k, v := range es {
		if k == nil || v == nil {
			continue
		}
		l.AddEntry(k, v)
	}
	return l
}

func (l *MapLiteral) asthLiteralExpr() ast.Expr { return l.expr }
func (l *MapLiteral) asthRValue() ast.Expr      { return l.expr }

// STRUCT

func NewStructLiteral() *StructLiteral {
	return &StructLiteral{&ast.CompositeLit{}}
}
func NewStructTypedLit(typ Type) *StructLiteral {
	return &StructLiteral{
		&ast.CompositeLit{
			Type: typ.asthType(),
		},
	}
}
func (l *StructLiteral) addField(f *StructFieldValue) *StructLiteral {
	l.expr.Elts = append(l.expr.Elts, f.expr)
	return l
}
func (l *StructLiteral) AddFields(fs ...*StructFieldValue) *StructLiteral {
	for _, f := range fs {
		if f == nil {
			continue
		}
		l.addField(f)
	}
	return l
}

func (l *StructLiteral) asthLiteralExpr() ast.Expr { return l.expr }
func (l *StructLiteral) asthRValue() ast.Expr      { return l.expr }

type StructFieldValue struct {
	expr ast.Expr
}

/// When definining values in a struct awhen specifying the field name
func NewStructFieldNamedVal(name string, val Rvalue) *StructFieldValue {
	return &StructFieldValue{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(name),
			Value: val.asthRValue(),
		},
	}
}

/// When definining values in a struct without specifying the field name
func NewStructFieldValue(lit Literal) *StructFieldValue {
	return &StructFieldValue{
		lit.asthLiteralExpr(),
	}
}
