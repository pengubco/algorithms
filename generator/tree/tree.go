// tree package defines methods that generates tree of specified shape.
package tree

import "fmt"

type Tree interface {
	Root() *Vertex
	Size() int
}

// FullTree is a tree where each non-leaf vertex has N children. N is the degree.
// For a full tree of height H and degree D, the total number of vertices is `(D^H -1)/(D-1)`.
type FullTree struct {
	root   *Vertex
	height int
	degree int
	size   int
}

// Vertex is a general vertex in a tree. Each vertex records information of its
// positional information in a tree.
// For simplicity, a leaf has empty Children[], instead a non-empty Children[]
// where each item in Children[] being nil.
type Vertex struct {
	// BFS order of the node in the tree. starting from 0.
	Order int

	// Level of the node in the tree. Root has level 0.
	Level int

	Children []*Vertex
	Parent   *Vertex

	// The next node in BFS order.
	Next *Vertex
}

func (v *Vertex) IsLeaf() bool {
	return len(v.Children) == 0
}

func (v *Vertex) IsRoot() bool {
	return v.Parent == nil
}

func (v *Vertex) IsInternal() bool {
	return !v.IsRoot() && !v.IsInternal()
}

// NewFullTree creates a full tree of height and degree.
func NewFullTree(height int, degree int) (*FullTree, error) {
	if height <= 0 || degree <= 0 {
		return nil, fmt.Errorf("invalid height or degree, height: %d, degree: %d",
			height, degree)
	}
	order := 0
	level := 0

	root := Vertex{
		Order:  order,
		Level:  level,
		Parent: nil,
		Next:   nil,
	}
	var previous *Vertex
	q := []*Vertex{&root}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if previous != nil {
			previous.Next = cur
		}
		previous = cur
		if cur.Level == height-1 {
			continue
		}
		cur.Children = make([]*Vertex, degree)
		for i := 0; i < degree; i++ {
			order++
			cur.Children[i] = &Vertex{
				Order:  order,
				Level:  cur.Level + 1,
				Parent: cur,
			}
			q = append(q, cur.Children[i])
		}
	}
	return &FullTree{
		root:   &root,
		height: height,
		degree: degree,
		size:   order + 1,
	}, nil
}

func (t *FullTree) Root() *Vertex {
	return t.root
}

func (t *FullTree) Size() int {
	return t.size
}

func (t *FullTree) Height() int {
	return t.height
}

func (t *FullTree) Degree() int {
	return t.degree
}
