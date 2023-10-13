package tree

import (
	"fmt"
	"math"
	"testing"

	"slices"

	"github.com/stretchr/testify/assert"
)

func TestFullTree(t *testing.T) {
	cases := []struct {
		height int
		degree int
	}{
		{1, 1},
		{1, 2},
		{2, 3},
		{10, 2},
	}

	// Test the structure of the full tree:
	// 1. DFS all vertices and sort vertices by their Order.
	// 2. Iterate all vertices, verify Order of its parent and children.
	for _, tc := range cases {
		t.Run(fmt.Sprintf("height: %d, degree %d", tc.height, tc.degree), func(t *testing.T) {
			assert := assert.New(t)
			ft, err := NewFullTree(tc.height, tc.degree)
			assert.Nil(err)
			vertices := getAllVerticesOrdered(ft.Root())
			treeSize := 0
			if tc.degree == 1 {
				treeSize = tc.height
			} else {
				treeSize = int(math.Pow(float64(tc.degree), float64(tc.height))-1) / (tc.degree - 1)
			}
			assert.Equal(treeSize, len(vertices))
			assert.Equal(treeSize, ft.Size())
			for i, v := range vertices {
				if i == 0 {
					assert.Equal(v, ft.Root())
				}
				if !v.IsRoot() {
					assert.Equal((v.Order-1)/tc.degree, v.Parent.Order)
				}
				if v.IsLeaf() {
					assert.Nil(v.Children)
					continue
				}
				assert.Equal(tc.degree, len(v.Children))
				for j := 0; j < tc.degree; j++ {
					assert.Equal(v.Order*tc.degree+j+1, v.Children[j].Order)
				}
			}
		})
	}
}

func getAllVerticesOrdered(root *Vertex) []*Vertex {
	var vertices []*Vertex
	var dfs func(v *Vertex)
	dfs = func(v *Vertex) {
		if v == nil {
			return
		}
		vertices = append(vertices, v)
		for _, c := range v.Children {
			dfs(c)
		}
	}
	dfs(root)
	slices.SortFunc(vertices, func(a, b *Vertex) int {
		return a.Order - b.Order
	})
	return vertices
}
