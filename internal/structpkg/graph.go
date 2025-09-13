package structpkg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

type Graph struct {
	adj map[string]Set
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[string]Set)}
}

func validateGraph(g *Graph) error {
	for vertex, neighbors := range g.adj {
		if vertex == "" {
			return fmt.Errorf("graph validation failed: vertex cannot be empty")
		}

		if strings.Contains(vertex, ";") {
			return fmt.Errorf("graph validation failed: vertex %q contains invalid character \";\"", vertex)
		}

		for neighbor := range neighbors {
			if _, ok := g.adj[neighbor]; !ok {
				return fmt.Errorf("graph validation failed: vertex %q referenced from %q does not exist", neighbor, vertex)
			}

			if _, ok := g.adj[neighbor][vertex]; !ok {
				return fmt.Errorf("graph validation failed: edge %q → %q is not symmetric", vertex, neighbor)
			}

			if neighbor == vertex {
				return fmt.Errorf("graph validation failed: self-loop detected at vertex %q", vertex)
			}
		}
	}

	return nil
}

func (g *Graph) bfs(start string, target ...string) Set {
	visited := make(Set)

	if !g.HasVertex(start) {
		return visited
	}

	visited[start] = Void{}
	queue := []string{start}
	current_index := 0
	targetSpecified := false

	if len(target) > 0 {
		targetSpecified = true
	}

	for len(queue) > current_index {
		current := queue[current_index]
		current_index++

		for neighbor := range g.adj[current] {
			if _, ok := visited[neighbor]; !ok {
				visited[neighbor] = Void{}

				if targetSpecified && neighbor == target[0] {
					return visited
				}

				queue = append(queue, neighbor)
			}
		}
	}

	return visited
}

func (g *Graph) AddVertex(vertex string) error {
	if g.HasVertex(vertex) {
		return nil
	}

	if strings.Contains(vertex, ";") {
		return fmt.Errorf("graph validation error: vertex %q contains invalid character \";\"", vertex)
	}

	if vertex == "" {
		return fmt.Errorf("graph validation error: vertex cannot be empty string")
	}

	g.adj[vertex] = make(Set)

	return nil
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

func (g *Graph) RemoveVertexIfIsolated(vertex string) {
	if !g.HasVertex(vertex) {
		return
	}

	if len(g.adj[vertex]) == 0 {
		delete(g.adj, vertex)
	}
}

func (g *Graph) GetVertices() []string {
	var vertices []string

	for vertex := range g.adj {
		vertices = append(vertices, vertex)
	}

	return vertices
}

func (g *Graph) GetConnectivityGroups() [][]string {
	visited := make(Set)
	var groups [][]string

	for vertex := range g.adj {
		if _, ok := visited[vertex]; ok {
			continue
		}

		var group []string

		for v := range g.bfs(vertex) {
			group = append(group, v)
			visited[v] = Void{}
		}

		groups = append(groups, group)
	}

	return groups
}

func (g *Graph) AddEdge(a, b string) error {
	if g.HasEdge(a, b) {
		return nil
	}

	if !g.HasVertex(a) {
		err := g.AddVertex(a)

		if err != nil {
			return err
		}
	}

	if a == b {
		return nil
	}

	if !g.HasVertex(b) {
		err := g.AddVertex(b)

		if err != nil {
			return err
		}
	}

	g.adj[a][b] = Void{}
	g.adj[b][a] = Void{}

	return nil
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

	g.RemoveVertexIfIsolated(a)
	g.RemoveVertexIfIsolated(b)
}

func (g *Graph) GetConnectedVertices(vertex string) []string {
	var connected []string

	for v := range g.bfs(vertex) {
		if vertex != v {
			connected = append(connected, v)
		}
	}

	return connected
}

func (g *Graph) ConnectedVertexCount(vertex string) int {
	visited := g.bfs(vertex)

	if len(visited) == 0 {
		return 0
	} else {
		return len(visited) - 1
	}
}

func (g *Graph) AreConnected(a, b string) bool {
	if !g.HasVertex(a) {
		return false
	}

	_, ok := g.bfs(a, b)[b]

	return ok
}

func (g *Graph) Cleanup() {
	for vertex := range g.adj {
		g.RemoveVertexIfIsolated(vertex)
	}
}

func (g *Graph) Order() int {
	return len(g.adj)
}

func (g *Graph) Clone() *Graph {
	clone := NewGraph()

	for vertex, neighbors := range g.adj {
		clonedNeighbors := make(Set)

		for neighbor := range neighbors {
			clonedNeighbors[neighbor] = Void{}
		}

		clone.adj[vertex] = clonedNeighbors
	}

	return clone
}

func (g *Graph) SerializeGob() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(g)

	return buf.Bytes()
}

func DeserializeGob(data []byte) (*Graph, error) {
	g := NewGraph()
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&g); err != nil {
		return nil, fmt.Errorf("graph deserialization failed: %w", err)
	}

	err := validateGraph(g)

	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Graph) SerializeCsv() []byte {
	var buf bytes.Buffer

	for vertex, neighbors := range g.adj {
		fmt.Fprint(&buf, vertex)

		for neighbor := range neighbors {
			fmt.Fprintf(&buf, ";%s", neighbor)
		}

		fmt.Fprint(&buf, '\n')
	}

	return buf.Bytes()
}

func DeserializeCsv(data []byte) (*Graph, error) {
	g := NewGraph()
	text := strings.TrimSpace(string(data))
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		vertices := strings.Split(line, ";")

		if _, ok := g.adj[vertices[0]]; ok {
			return nil, fmt.Errorf("graph validation failed: duplicate vertex %q found", vertices[0])
		}

		g.adj[vertices[0]] = make(Set)

		for i := 1; i < len(vertices); i++ {
			g.adj[vertices[0]][vertices[i]] = Void{}
		}
	}

	err := validateGraph(g)

	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Graph) SerializeCsvCondensed() []byte {
	var buf bytes.Buffer

	for vertex, neighbors := range g.adj {
		if len(neighbors) == 0 {
			fmt.Fprintf(&buf, "%s\n", vertex)

			continue
		}

		for neighbor := range neighbors {
			if neighbor > vertex {
				fmt.Fprintf(&buf, "%s;%s\n", vertex, neighbor)
			}
		}
	}

	return buf.Bytes()
}

func DeserializeCsvCondensed(data []byte) (*Graph, error) {
	g := NewGraph()
	text := strings.TrimSpace(string(data))
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		vertices := strings.Split(line, ";")
		var err error

		switch len(vertices) {
		case 1:
			err = g.AddVertex(vertices[0])

		case 2:
			err = g.AddEdge(vertices[0], vertices[1])

		default:
			err = fmt.Errorf("graph deserialization failed: invalid number of vertices in line %q", line)
		}

		if err != nil {
			return nil, err
		}
	}

	return g, nil
}
