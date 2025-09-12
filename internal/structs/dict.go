package structs

type Dict struct {
	graph *Graph
}

func NewDict() *Dict {
	return &Dict{graph: NewGraph()}
}

func (d *Dict) AddSynonims(a, b string) {
	d.graph.AddEdge(a, b)
}

func (d *Dict) AddWord(word string) {
	d.graph.AddVertex(word)
}

func (d *Dict) RemoveWord(words ...string) {
	for _, word := range words {
		d.graph.RemoveVertex(word)
	}
}

func (d *Dict) UnlinkSynonims(a, b string) {
	d.graph.RemoveEdge(a, b)
}

func (d *Dict) GetSymomims(word string) {
	d.graph.GetConnectedVertices(word)
}

func (d *Dict) SynonimCount(word string) {
	d.graph.ConnectedVertexSize(word)
}
