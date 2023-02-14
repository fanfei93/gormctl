package genstruct

type GenElement struct {
	Name       string
	ColumnName string
	Type       string
	Notes      string
	Tags       map[string][]string
}

type GenStruct struct {
	Name string
	Em   []GenElement
}

type GenPackage struct {
	Name        string
	Imports     map[string]string
	Structs     []GenStruct
	FuncStrList []string
}
