package main

import "fmt"

func (g *digraph) shortestPath(src, dest string) (string, error) {
	start, end := g.vertices[src], g.vertices[dest]

	tree := bfs(start)

	if prev, ok := tree[end]; ok {
		result := end.element

		for prev != nil {
			result = prev.element + " -> " + result
			prev = tree[prev]
		}
		return result, nil
	}
	return "", fmt.Errorf("destination not reachable from source")
}

func bfs(src *vertex) map[*vertex]*vertex {
	tree := map[*vertex]*vertex{}
	visited := map[*vertex]bool{}

	q := []*vertex{}

	visited[src] = true
	q = append(q, src)

	for len(q) != 0 {
		p := q[0]
		q = q[1:]

		for _, v := range p.out {
			if t := visited[v]; !t {
				tree[v] = p
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	return tree
}
