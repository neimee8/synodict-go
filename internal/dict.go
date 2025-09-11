package internal

type Dict struct {
	graph *Graph
}

func NewDict() *Dict {
	return &Dict{graph: NewGraph()}
}
