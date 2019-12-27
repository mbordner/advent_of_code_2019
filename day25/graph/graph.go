package graph

type Edge struct {
	source      *Node
	destination *Node
	value       float64
	properties  map[string]interface{}
}

func (e *Edge) IsTraversable() bool {
	return e.destination.IsTraversable()
}

func (e *Edge) GetSource() *Node {
	return e.source
}

func (e *Edge) GetDestination() *Node {
	return e.destination
}

func (e *Edge) SetDestination(o *Node) {
	e.destination = o
}

func (e *Edge) GetValue() float64 {
	return e.value
}

func (e *Edge) AddProperty(id string, value interface{}) {
	e.properties[id] = value
}

func (e *Edge) GetProperty(id string) interface{} {
	if v, ok := e.properties[id]; ok {
		return v
	}
	return nil
}

type Node struct {
	id          interface{}
	edges       []*Edge
	properties  map[string]interface{}
	traversable bool
}

func (n *Node) GetID() interface{} {
	return n.id
}

func (n *Node) GetEdges() []*Edge {
	edges := make([]*Edge, len(n.edges), len(n.edges))
	for i := range n.edges {
		edges[i] = n.edges[i]
	}
	return edges
}

func (n *Node) GetTraversableEdges() []*Edge {
	edges := make([]*Edge, 0, len(n.edges))
	for i := range n.edges {
		if n.edges[i].IsTraversable() {
			edges = append(edges, n.edges[i])
		}
	}
	return edges
}

func (n *Node) IsTraversable() bool {
	return n.traversable
}

func (n *Node) AddProperty(id string, value interface{}) {
	n.properties[id] = value
}

func (n *Node) GetProperty(id string) interface{} {
	if v, ok := n.properties[id]; ok {
		return v
	}
	return nil
}

func (n *Node) SetTraversable(b bool) {
	n.traversable = b
}

func (n *Node) AddEdge(o *Node, w float64) *Edge {
	e := Edge{source: n, destination: o, value: w}
	e.properties = make(map[string]interface{})
	if n.edges == nil {
		n.edges = make([]*Edge, 0, 8)
	}
	n.edges = append(n.edges, &e)
	return &e
}

type Graph struct {
	nodes map[interface{}]*Node
}

func NewGraph() *Graph {
	g := new(Graph)
	g.nodes = make(map[interface{}]*Node)
	return g
}

func (g *Graph) Len() int {
	return len(g.nodes)
}

func (g *Graph) CreateNode(id interface{}) *Node {
	n := new(Node)
	n.id = id
	n.properties = make(map[string]interface{})
	n.traversable = true
	g.nodes[n.id] = n
	return n
}

func (g *Graph) GetNode(id interface{}) *Node {
	if n, ok := g.nodes[id]; ok {
		return n
	}
	return nil
}

func (g *Graph) GetNodes() []*Node {
	ns := make([]*Node, len(g.nodes), len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		ns[i] = n
		i++
	}
	return ns
}

func (g *Graph) GetTraversableNodes() []*Node {
	ns := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		if n.IsTraversable() {
			ns = append(ns, n)
		}
	}
	return ns
}
