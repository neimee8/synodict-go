package structpkg

type Dict struct {
	graph *Graph
}

func NewDict() *Dict {
	return &Dict{graph: NewGraph()}
}

func (d *Dict) AddSynonyms(words ...string) {
	for i := 0; i < len(words)-1; i++ {
		d.graph.AddEdge(words[i], words[i+1])
	}
}

func (d *Dict) AddWords(words ...string) {
	for _, word := range words {
		d.graph.AddVertex(word)
	}
}

func (d *Dict) RemoveWords(words ...string) {
	for _, word := range words {
		d.graph.RemoveVertex(word)
	}
}

func (d *Dict) UnlinkSynonyms(a, b string) {
	d.graph.RemoveEdge(a, b)
}

func (d *Dict) UnlinkSynonymsAndCleanup(a, b string) {
	d.graph.RemoveEdgeAndCleanup(a, b)
}

func (d *Dict) GetSymomims(word string) []string {
	return d.graph.GetConnectedVertices(word)
}

func (d *Dict) SynonymCount(word string) int {
	return d.graph.ConnectedVertexCount(word)
}

func (d *Dict) GetWords() []string {
	return d.graph.GetVertices()
}

func (d *Dict) GetSynonymGroups() [][]string {
	return d.graph.GetConnectivityGroups()
}

func (d *Dict) Clear() {
	d.graph = NewGraph()
}

func (d *Dict) Cleanup() {
	d.graph.Cleanup()
}
