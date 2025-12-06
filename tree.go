package listmd

import (
	"sort"
	"strconv"
)

func buildTree(Nodes []Node) [][]*Node {
	groupedNodes := make(map[int][]Node)
	for _, node := range Nodes {
		groupedNodes[node.collectionIndex] = append(groupedNodes[node.collectionIndex], node)
	}

	var collectionIndices []int
	for idx := range groupedNodes {
		collectionIndices = append(collectionIndices, idx)
	}
	sort.Ints(collectionIndices)

	var result [][]*Node
	for _, collectionIdx := range collectionIndices {
		nodes := groupedNodes[collectionIdx]
		roots := buildTreeForCollection(nodes, collectionIdx)
		result = append(result, roots)
	}

	return result
}

func buildTreeForCollection(Nodes []Node, collectionIdx int) []*Node {
	// Create a map to quickly look up nodes by ID
	nodeMap := make(map[int]*Node)

	// Sort flat nodes by ID first
	sort.Slice(Nodes, func(i, j int) bool {
		return Nodes[i].Id < Nodes[j].Id
	})

	// First pass: create all nodes without children
	for _, n := range Nodes {
		node := &Node{
			Id:       n.Id,
			parentId: n.parentId,
			level:    n.level,
			Markdown: n.Markdown,
			Children: []*Node{},
		}
		nodeMap[n.Id] = node
	}

	// Second pass: build the tree structure
	var roots []*Node
	for _, node := range nodeMap {
		if node.parentId == -1 {
			// This is a root node
			roots = append(roots, node)
		} else {
			// Find parent and add this node as a child
			if parent, exists := nodeMap[node.parentId]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	// Sort roots by ID
	sortNodesByID(roots)

	// Sort children recursively
	for _, root := range roots {
		sortTree(root)
	}

	// Assign UIDs starting from the roots
	assignUIDs(roots, collectionIdx)

	return roots
}

func sortTree(node *Node) {
	if len(node.Children) > 0 {
		sortNodesByID(node.Children)
		for _, child := range node.Children {
			sortTree(child)
		}
	}
}

// sortNodesByID sorts a slice of nodes by their ID
func sortNodesByID(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id < nodes[j].Id
	})
}

func assignUIDs(nodes []*Node, collectionIdx int) {
	for i, node := range nodes {
		node.UID = "item_" + strconv.Itoa(collectionIdx) + "_" + strconv.Itoa(i)
		if len(node.Children) > 0 {
			assignChildUIDs(node.Children, node.UID)
		}
	}
}

func assignChildUIDs(children []*Node, parentUID string) {
	for i, child := range children {
		child.UID = parentUID + "_" + strconv.Itoa(i)
		if len(child.Children) > 0 {
			assignChildUIDs(child.Children, child.UID)
		}
	}
}
