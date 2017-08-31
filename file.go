package asth

import (
	"go/ast"
	"go/token"
	"strings"
)

type File struct {
	node *ast.File

	importDecl *ast.GenDecl
	importDoc  *ast.CommentGroup

	objects map[string]*ast.Object
}

func NewFile(pack string) *File {
	file := File{
		node: &ast.File{
			Name: ast.NewIdent(pack),
		},
		objects: map[string]*ast.Object{},
	}
	return &file
}

func (f *File) DefinedObject(name string) Lvalue {
	o, ok := f.objects[name]
	if !ok {
		return nil
	}

	return &BaseLvalue{expr: ast.NewIdent(o.Name)}
}

func (f *File) defineObject(obj *ast.Object) {
	f.objects[obj.Name] = obj
}

func (f *File) WithImportDoc(lines ...string) *File {
	list := []*ast.Comment{}

	for _, l := range lines {
		l = "///" + l
		list = append(list, &ast.Comment{Text: l})
	}
	f.importDoc = &ast.CommentGroup{List: list}
	if f.importDecl != nil {
		f.importDecl.Doc = f.importDoc
	}
	return f
}

func (f *File) Get() *ast.File {
	return f.node
}

func (f *File) AddImport(imp *ImportSpec) {
	if f.importDecl == nil {
		f.importDecl = &ast.GenDecl{
			Tok:    token.IMPORT,
			Specs:  []ast.Spec{},
			Lparen: 1, // We just need something valid (!=0)
			Rparen: 1, // We just need something valid (!=0)
		}
		if f.importDoc != nil {
			f.importDecl.Doc = f.importDoc
		}
		f.node.Decls = append([]ast.Decl{f.importDecl}, f.node.Decls...)
	}
	f.importDecl.Specs = append(f.importDecl.Specs, imp.spec)
	f.node.Imports = append(f.node.Imports, imp.spec)
}
func (f *File) AddImports(imp ...*ImportSpec) {
	for _, i := range imp {
		f.AddImport(i)
	}
	ast.SortImports(token.NewFileSet(), f.node)
}

func (f *File) AddDecl(decl Decl) {
	switch d := decl.(type) {
	case *GenDecl:
		switch d.node.Tok {
		case token.VAR:
			for _, v := range d.node.Specs {
				vspec, ok := v.(*ast.ValueSpec)
				if !ok {
					panic("VAR GenDecl contains a non ValueSpec spec.")
				}
				for _, name := range vspec.Names {
					if strings.Contains(name.Name, ".") {
						// Ignore name containing a .,  they're not defined in the file scope
						continue
					}
					f.defineObject(name.Obj)
				}
			}
		}
	}

	f.node.Decls = append(f.node.Decls, decl.asthDeclNode())
}
