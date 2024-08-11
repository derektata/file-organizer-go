package main

import (
	"fmt"
	"os"
	"strings"
)

// Directory structure map to hold the tree view
type DirectoryTree struct {
	Root *DirectoryNode
}

// DirectoryNode represents a node in the directory tree
type DirectoryNode struct {
	Name     string
	Children map[string]*DirectoryNode
	IsFile   bool
}

// AddFile adds a file to the directory tree structure.
//
// It takes a filePath as a parameter and updates the DirectoryTree accordingly.
// The filePath is a string representing the path of the file to be added.
// Return type: None.
func (tree *DirectoryTree) AddFileToTree(filePath string) {
	relativePath := strings.TrimPrefix(filePath, tree.Root.Name)
	parts := strings.Split(relativePath, string(os.PathSeparator))
	current := tree.Root

	for i, part := range parts {
		if part == "" {
			continue
		}
		if _, exists := current.Children[part]; !exists {
			current.Children[part] = &DirectoryNode{
				Name:     part,
				Children: make(map[string]*DirectoryNode),
				IsFile:   i == len(parts)-1,
			}
		}
		current = current.Children[part]
	}
}

// PrintSubTree prints the subtree of the directory tree rooted at the specified path.
//
// Parameter rootPath is the path of the directory from which to start printing the subtree.
// Return type: None.
func (tree *DirectoryTree) PrintSubTree(rootPath string) {
	node := tree.findNode(rootPath)
	if node == nil {
		fmt.Printf("Directory %s not found in the tree.\n", rootPath)
		return
	}

	fmt.Println(node.Name)
	for _, child := range node.Children {
		switch {
		case child.IsFile:
			fmt.Printf("    └── %s\n", child.Name)
		case len(child.Children) > 0:
			fmt.Printf("    ├── %s\n", child.Name)
			for _, grandChild := range child.Children {
				fmt.Printf("    │   └── %s\n", grandChild.Name)
			}
		default:
			fmt.Printf("    ├── %s\n", child.Name)
		}
	}
}

// findNode finds the DirectoryNode corresponding to the specified path in the DirectoryTree.
//
// Parameter path is the absolute path of the directory or file to be found.
// Return type is a pointer to the DirectoryNode if found, otherwise nil.
func (tree *DirectoryTree) findNode(path string) *DirectoryNode {
	relativePath := strings.TrimPrefix(path, tree.Root.Name)
	parts := strings.Split(relativePath, string(os.PathSeparator))
	current := tree.Root

	for _, part := range parts {
		if part == "" {
			continue
		}
		if child, exists := current.Children[part]; exists {
			current = child
		} else {
			return nil
		}
	}

	return current
}
