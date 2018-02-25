package lattice

func Do(run func(int, DataType) (NodeIterator, error), support int, dt DataType, do func(Node) error) error {
	it, err := run(support, dt)
	if err != nil {
		return err
	}
	var node Node
	for node, err, it = it(); it != nil; node, err, it = it() {
		e := do(node)
		if e != nil {
			return e
		}
	}
	return err
}

func Slice(run func(int, DataType) (NodeIterator, error), support int, dt DataType) ([]Node, error) {
	nodes := make([]Node, 0, 10)
	err := Do(run, support, dt,
		func(n Node) error {
			nodes = append(nodes, n)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func NodeIteratorFromSlice(nodes []Node) (it NodeIterator, err error) {
	i := 0
	it = func() (Node, error, NodeIterator) {
		if i >= len(nodes) {
			return nil, nil, nil
		}
		n := nodes[i]
		i++
		return n, nil, it
	}
	return it, nil
}
