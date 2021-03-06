package asth

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

type (
	Package struct {
		name  *ast.Ident
		files map[string]*File
	}
)

func NewPackage(name string) *Package {
	return &Package{
		name:  ast.NewIdent(name),
		files: map[string]*File{},
	}
}

func (p *Package) NewFile(name string) *File {
	file := &File{
		Name: name,
		node: &ast.File{
			Name: p.name,
		},
		objects: map[string]*ast.Object{},
	}
	p.files[name] = file
	return file
}

func (p *Package) DefinedObject(name string) Lvalue {
	for _, f := range p.files {
		if o := f.DefinedObject(name); o != nil {
			return o
		}
	}

	defined := p.ListObjects()
	err := fmt.Errorf("Unknown object `%s`. Defined objects are:\n\t%s", name, strings.Join(defined, "\n\t"))
	log.Println("Warning:", err)
	return nil
}

func (p *Package) ListObjects() []string {
	objs := []string{}

	for _, f := range p.files {
		objs = append(f.ListObjects(), objs...)
	}
	return objs
}

func (p *Package) WriteFiles(outDir string, cfg printer.Config) error {
	fset := token.NewFileSet()
	files := p.files

	for k, v := range files {
		log.Println("Exporting file", k)
		path := fmt.Sprintf("%s/%s", outDir, k)
		out, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error opening %s file: %s", k, err)
		}
		if err := cfg.Fprint(out, fset, v.Get()); err != nil {
			return fmt.Errorf("Error writing %s file: %s", k, err)
		}
	}
	return nil
}
