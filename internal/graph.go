package internal

type Graph struct {
	adj map[string]Set
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[string]Set)}
}

func (g *Graph) bfs(start string) map[string]bool {
	visited := make(map[string]bool)

	if !g.HasVertex(start) {
		return visited
	}

	visited[start] = true
	queue := []string{start}
	current_index := 0

	for len(queue) > current_index {
		current := queue[current_index]
		current_index++

		for neighbor := range g.adj[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return visited
}

func (g *Graph) AddVertex(vertex string) {
	if g.HasVertex(vertex) {
		return
	}

	g.adj[vertex] = Set{}
}

func (g *Graph) HasVertex(vertex string) bool {
	_, ok := g.adj[vertex]
	return ok
}

func (g *Graph) RemoveVertex(vertex string) {
	if !g.HasVertex(vertex) {
		return
	}

	neighbors := g.adj[vertex]

	for neighbor := range neighbors {
		delete(g.adj[neighbor], vertex)
	}

	delete(g.adj, vertex)
}

func (g *Graph) AddEdge(a, b string) {
	if g.HasEdge(a, b) {
		return
	}

	if !g.HasVertex(a) {
		g.AddVertex(a)
	}

	if a == b {
		return
	}

	if !g.HasVertex(b) {
		g.AddVertex(b)
	}

	g.adj[a][b] = Void{}
	g.adj[b][a] = Void{}
}

func (g *Graph) HasEdge(a, b string) bool {
	if !g.HasVertex(a) {
		return false
	}

	_, ok := g.adj[a][b]
	return ok
}

func (g *Graph) RemoveEdge(a, b string) {
	if !g.HasEdge(a, b) {
		return
	}

	delete(g.adj[a], b)
	delete(g.adj[b], a)
}

func (g *Graph) RemoveEdgeAndCleanup(a, b string) {
	g.RemoveEdge(a, b)

	if len(g.adj[a]) == 0 {
		g.RemoveVertex(a)
	}

	if len(g.adj[b]) == 0 {
		g.RemoveVertex(b)
	}
}

func (g *Graph) GetConnectedVertexes(vertex string) []string {
	var connected []string

	for v := range g.bfs(vertex) {
		connected = append(connected, v)
	}

	return connected
}

func (g *Graph) ConnectedVertexSize(vertex string) int {
	return len(g.GetConnectedVertexes(vertex))
}

func (g *Graph) AreConnected(a, b string) bool {
	if !g.HasVertex(a) {
		return false
	}

	return g.bfs(a)[b]
}

func (g *Graph) Cleanup() {
	for vertex, edges := range g.adj {
		if len(edges) == 0 {
			delete(g.adj, vertex)
		}
	}
}

func (g *Graph) Order() int {
	return len(g.adj)
}
