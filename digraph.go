package main

import (
	"fmt"
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

func (v *vertex) customString() string {
	var out, in string

	if len(v.in) > 0 {
		in = v.in[0].element

		for _, r := range v.in[1:] {
			in += "," + r.element
		}
	}

	if len(v.out) > 0 {
		out = v.out[0].element

		for _, r := range v.out[1:] {
			out += "," + r.element
		}
	}

	return fmt.Sprintf("%v;%v;%v", v.element, in, out)
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

func (g *digraph) getVertex(d string) *vertex {
	g.lock.RLock()
	v, _ := g.vertices[d]
	g.lock.RUnlock()
	return v
}
