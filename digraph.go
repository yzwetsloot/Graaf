package main

import (
	"fmt"
	"os"
	"sync"
)

type vertex struct {
	element string
	in      []*vertex
	out     []*vertex
	lock    sync.RWMutex
}

func (v *vertex) addIncoming(u *vertex) {
	v.lock.Lock()
	v.in = append(v.in, u)
	v.lock.Unlock()
}

func (v *vertex) addOutgoing(u *vertex) {
	v.lock.Lock()
	v.out = append(v.out, u)
	v.lock.Unlock()
}

func (v *vertex) String() string {
	return fmt.Sprintf("element: %v, in-degree: %v, out-degree: %v", v.element, len(v.in), len(v.out))
}

func (v *vertex) expandString() string {
	var out, in string

	if len(v.in) > 0 {
		in = v.in[0].element

		for _, u := range v.in[1:] {
			in += "," + u.element
		}
	}

	if len(v.out) > 0 {
		out = v.out[0].element

		for _, u := range v.out[1:] {
			out += "," + u.element
		}
	}

	return fmt.Sprintf("%v;%v;%v", v.element, in, out)
}

type path []*vertex

func (p path) String() string {
	if len(p) > 0 {
		result := p[len(p)-1].element

		for _, v := range p[:len(p)-1] {
			result = v.element + " -> " + result
		}

		return result
	}
	return ""
}

type digraph struct {
	vertices map[string]*vertex
	lock     sync.RWMutex
}

func (g *digraph) addVertex(v *vertex) {
	g.lock.Lock()
	g.vertices[v.element] = v
	g.lock.Unlock()
}

func (g *digraph) containsDomain(d string) bool {
	g.lock.RLock()
	_, ok := g.vertices[d]
	g.lock.RUnlock()
	return ok
}

func (g *digraph) getVertex(d string) (*vertex, bool) {
	g.lock.RLock()
	v, ok := g.vertices[d]
	g.lock.RUnlock()
	return v, ok
}

func (g *digraph) String() string {
	outdeg := 0

	for v := range g.vertices {
		t, _ := g.getVertex(v)
		outdeg += len(t.out)
	}

	return fmt.Sprintf("number of nodes: %v, number of edges: %v", len(g.vertices), outdeg)
}

func (g *digraph) expandString() (result string) {
	sep := "\n"

	for _, v := range g.vertices {
		result += v.expandString() + sep
	}

	return
}

func (g *digraph) serialize(path string) {
	file, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	file.WriteString(g.expandString())
}
