package graph

import (
	"fmt"
)

func (g *Digraph) ShortestPath(src, dest string) (Path, error) {
	start, end := g.Vertices[src], g.Vertices[dest]

	tree := bfs(start)

	var p Path

	if prev, ok := tree[end]; ok {
		p = append(Path{end}, p...)

		for prev != nil {
			p = append(Path{prev}, p...)
			prev = tree[prev]
		}
		return p, nil
	}
	return p, fmt.Errorf("destination not reachable from source")
}

func bfs(src *Vertex) map[*Vertex]*Vertex {
	tree := map[*Vertex]*Vertex{}
	visited := map[*Vertex]bool{}

	var q []*Vertex

	visited[src] = true
	q = append(q, src)

	for len(q) != 0 {
		p := q[0]
		q = q[1:]

		for _, v := range p.Out {
			if t := visited[v]; !t {
				tree[v] = p
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	return tree
}
