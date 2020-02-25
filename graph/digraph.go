package graph

import (
	"fmt"
	"os"
	"sync"
)

type Vertex struct {
	Element string
	In      []*Vertex
	Out     []*Vertex
	Lock    sync.RWMutex
}

func (v *Vertex) AddIncoming(u *Vertex) {
	v.Lock.Lock()
	v.In = append(v.In, u)
	v.Lock.Unlock()
}

func (v *Vertex) AddOutgoing(u *Vertex) {
	v.Lock.Lock()
	v.Out = append(v.Out, u)
	v.Lock.Unlock()
}

func (v *Vertex) String() string {
	return fmt.Sprintf("element: %v, in-degree: %v, out-degree: %v", v.Element, len(v.In), len(v.Out))
}

func (v *Vertex) expandString() string {
	var out, in string

	if len(v.In) > 0 {
		in = v.In[0].Element

		for _, u := range v.In[1:] {
			in += "," + u.Element
		}
	}

	if len(v.Out) > 0 {
		out = v.Out[0].Element

		for _, u := range v.Out[1:] {
			out += "," + u.Element
		}
	}

	return fmt.Sprintf("%v;%v;%v", v.Element, in, out)
}

type Path []*Vertex

func (p Path) String() string {
	if len(p) > 0 {
		result := p[len(p)-1].Element

		for _, v := range p[:len(p)-1] {
			result = v.Element + " -> " + result
		}

		return result
	}
	return ""
}

type Digraph struct {
	Vertices map[string]*Vertex
	Lock     sync.RWMutex
}

func (g *Digraph) AddVertex(v *Vertex) {
	g.Lock.Lock()
	g.Vertices[v.Element] = v
	g.Lock.Unlock()
}

func (g *Digraph) containsDomain(d string) bool {
	g.Lock.RLock()
	_, ok := g.Vertices[d]
	g.Lock.RUnlock()
	return ok
}

func (g *Digraph) GetVertex(d string) (*Vertex, bool) {
	g.Lock.RLock()
	v, ok := g.Vertices[d]
	g.Lock.RUnlock()
	return v, ok
}

func (g *Digraph) String() string {
	outdeg := 0

	for v := range g.Vertices {
		t, _ := g.GetVertex(v)
		outdeg += len(t.Out)
	}

	return fmt.Sprintf("number of nodes: %v, number of edges: %v", len(g.Vertices), outdeg)
}

func (g *Digraph) expandString() (result string) {
	sep := "\n"

	for _, v := range g.Vertices {
		result += v.expandString() + sep
	}

	return
}

func (g *Digraph) Serialize(path string) {
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	file.WriteString(g.expandString())
}
