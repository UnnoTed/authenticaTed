package generator

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func Parse(filePath string, data []byte) (*Skeleton, error) {
	var list = &Skeleton{
		Fields: make(map[string]*Field),
	}

	/*var conf loader.Config
	f, err := conf.ParseFile(filePath, data)
	if err != nil {
		return err
	}

	conf.CreateFromFiles("authenticaTed", f)
	prog, err := conf.Load()
	if err != nil {
		return err
	}

	//spew.Dump(prog.Created[0].Types)
	for _, d := range prog.Created[0].Types {
		t := reflect.ValueOf(d.Type).Elem()
		if d.IsType() && t.Type().String() == "types.Struct" {
			ty := structs.Map(d)["Type"]
			a :=
		}
	}*/

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, 0)
	if err != nil {
		panic(err)
	}

	var found bool
	ast.Inspect(f, func(n ast.Node) bool {
		if t, ok := n.(*ast.TypeSpec); ok {
			// get the struct's info
			s := t.Type.(*ast.StructType)
			fields := s.Fields.List

			for _, field := range fields {
				pos := field.Type.Pos() - 1
				end := field.Type.End() - 1
				ts := data[pos:end]

				var tags string
				if field.Tag != nil {
					tags = string(data[field.Tag.Pos() : field.Tag.End()-2])
				}

				name := field.Names[0].String()
				list.Fields[name] = &Field{
					Name: name,
					Type: string(ts),
					Tags: tags,
				}
			}

			found = true
			return false
		}

		return true
	})

	if !found {
		panic(errors.New("Error: no type found"))
	}

	log.Println(list)

	return list, nil
}
