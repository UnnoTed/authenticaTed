package generator

type Field struct {
	Name string
	Type string
	Tags string
}

type Skeleton struct {
	Fields   map[string]*Field
	Embedded *Skeleton
}
