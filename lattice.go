package lattice

import (
	"bytes"
)

func MakeLattice(n Node) (*Lattice, error) {
	lat, err := n.Lattice()
	if err != nil {
		_, ok := err.(*NoLattice)
		if !ok {
			return nil, err
		}
	} else {
		return lat, nil
	}
	return lattice(n)
}

func lattice(node Node) (*Lattice, error) {
	pop := func(queue []Node) (Node, []Node) {
		n := queue[0]
		copy(queue[0:len(queue)-1], queue[1:len(queue)])
		queue = queue[0 : len(queue)-1]
		return n, queue
	}
	queue := make([]Node, 0, 10)
	queue = append(queue, node)
	queued := make(map[string]bool)
	rlattice := make([]Node, 0, 10)
	for len(queue) > 0 {
		var n Node
		n, queue = pop(queue)
		nlabel := n.Pattern().Label()
		queued[string(nlabel)] = true
		rlattice = append(rlattice, n)
		parents, err := n.Parents()
		if err != nil {
			return nil, err
		}
		for _, p := range parents {
			haskid := false
			pkids, err := p.Children()
			if err != nil {
				return nil, err
			}
			for _, k := range pkids {
				if bytes.Equal(k.Pattern().Label(), nlabel) {
					haskid = true
					break
				}
			}
			if !haskid {
				// we were not able to compute the child from the parent
				// so let's drop this parent.
				continue
			}
			l := string(p.Pattern().Label())
			if _, has := queued[l]; !has {
				queue = append(queue, p)
				queued[l] = true
			}
		}
	}
	lattice := make([]Node, 0, len(rlattice))
	labels := make(map[string]int, len(lattice))
	for i := len(rlattice) - 1; i >= 0; i-- {
		lattice = append(lattice, rlattice[i])
		labels[string(lattice[len(lattice)-1].Pattern().Label())] = len(lattice) - 1
	}
	edges := make([]Edge, 0, len(lattice)*2)
	lattice_kids := make([][]*Edge, 0, len(lattice)*2)
	for i, n := range lattice {
		kids, err := n.Children()
		if err != nil {
			return nil, err
		}
		lattice_kids = append(lattice_kids, make([]*Edge, 0, len(kids)))
		for _, kid := range kids {
			j, has := labels[string(kid.Pattern().Label())]
			if has {
				edges = append(edges, Edge{Src: i, Targ: j})
				e := &edges[len(edges)-1]
				lattice_kids[len(lattice_kids)-1] = append(lattice_kids[len(lattice_kids)-1], e)
			}
		}
	}
	return &Lattice{lattice, edges, lattice_kids}, nil
}

func (lat *Lattice) Children(i int) []Node {
	kids := make([]Node, 0, len(lat.Kids[i]))
	for _, e := range lat.Kids[i] {
		kids = append(kids, lat.V[e.Targ])
	}
	return kids
}
