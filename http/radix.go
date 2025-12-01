package http

import (
	"strings"
)

// NodeKind represents the type of a radix tree node
type NodeKind uint8

const (
	NodeStatic   NodeKind = iota // Static segment (e.g., "users")
	NodeParam                    // Parameter segment (e.g., ":id")
	NodeWildcard                 // Wildcard segment (e.g., "*path")
)

// RadixNode represents a node in the radix tree
type RadixNode struct {
	prefix    string   // Path segment (e.g., "users", "api")
	kind      NodeKind // Type of node
	paramName string   // Parameter name for param/wildcard nodes

	// Separated children for fast dispatch:
	// - Static: slice for linear search (most routes have few children)
	// - Param: single pointer (only one :param per level allowed)
	// - Wildcard: single pointer (only one *wildcard per level allowed)
	staticChildren []*RadixNode
	paramChild     *RadixNode
	wildcardChild  *RadixNode

	// Handler for this node (nil if not terminal)
	handler        Handler
	chainedHandler Handler
}

// RadixTree is the main routing tree with method-indexed roots
type RadixTree struct {
	// Method-indexed trees for O(1) method dispatch
	// Index by Method value (MethodGet=0, MethodPost=1, etc.)
	trees [9]*RadixNode
}

// ParamCapture holds a captured path parameter
type ParamCapture struct {
	Name  string
	Value string
}

// NewRadixTree creates a new empty radix tree
func NewRadixTree() *RadixTree {
	return &RadixTree{}
}

// Insert adds a route to the radix tree
func (t *RadixTree) Insert(method Method, pattern string, handler, chainedHandler Handler) {
	// Get or create root for this method
	root := t.trees[method]
	if root == nil {
		root = &RadixNode{
			kind: NodeStatic,
		}
		t.trees[method] = root
	}

	// Handle root path specially
	if pattern == "/" || pattern == "" {
		root.handler = handler
		root.chainedHandler = chainedHandler
		return
	}

	// Parse pattern into segments
	segments := parsePatternSegments(pattern)

	// Walk tree, creating nodes as needed
	current := root
	for _, seg := range segments {
		current = insertOrFindChild(current, seg)
	}

	// Set handler on terminal node
	current.handler = handler
	current.chainedHandler = chainedHandler
}

// segmentInfo holds parsed segment data
type segmentInfo struct {
	kind  NodeKind
	value string // For static: the segment text; for param/wildcard: the param name
}

// parsePatternSegments parses a pattern into segments
func parsePatternSegments(pattern string) []segmentInfo {
	// Remove leading/trailing slashes
	pattern = strings.Trim(pattern, "/")
	if pattern == "" {
		return nil
	}

	parts := strings.Split(pattern, "/")
	segments := make([]segmentInfo, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		var seg segmentInfo
		if strings.HasPrefix(part, ":") {
			seg = segmentInfo{
				kind:  NodeParam,
				value: strings.TrimPrefix(part, ":"),
			}
		} else if strings.HasPrefix(part, "*") {
			seg = segmentInfo{
				kind:  NodeWildcard,
				value: strings.TrimPrefix(part, "*"),
			}
		} else {
			seg = segmentInfo{
				kind:  NodeStatic,
				value: part,
			}
		}
		segments = append(segments, seg)
	}

	return segments
}

// insertOrFindChild finds or creates a child node for the given segment
func insertOrFindChild(parent *RadixNode, seg segmentInfo) *RadixNode {
	switch seg.kind {
	case NodeStatic:
		// Look for existing static child with same prefix
		for _, child := range parent.staticChildren {
			if child.prefix == seg.value {
				return child
			}
		}
		// Create new static child
		newNode := &RadixNode{
			prefix: seg.value,
			kind:   NodeStatic,
		}
		parent.staticChildren = append(parent.staticChildren, newNode)
		return newNode

	case NodeParam:
		// Only one param child allowed per node
		if parent.paramChild == nil {
			parent.paramChild = &RadixNode{
				kind:      NodeParam,
				paramName: seg.value,
			}
		}
		return parent.paramChild

	case NodeWildcard:
		// Only one wildcard child allowed per node
		if parent.wildcardChild == nil {
			parent.wildcardChild = &RadixNode{
				kind:      NodeWildcard,
				paramName: seg.value,
			}
		}
		return parent.wildcardChild
	}

	return nil
}

// Lookup finds a handler for the given method and path
// Returns the matching node and captured parameters
func (t *RadixTree) Lookup(method Method, path string) (*RadixNode, []ParamCapture) {
	root := t.trees[method]
	if root == nil {
		return nil, nil
	}

	// Handle root path
	if path == "/" || path == "" {
		if root.chainedHandler != nil {
			return root, nil
		}
		return nil, nil
	}

	// Pre-allocate params slice (most routes have 0-4 params)
	params := make([]ParamCapture, 0, 4)

	node := t.lookupRecursive(root, path, 0, &params)
	if node != nil {
		return node, params
	}
	return nil, nil
}

// lookupRecursive performs recursive tree traversal
func (t *RadixTree) lookupRecursive(node *RadixNode, path string, pos int, params *[]ParamCapture) *RadixNode {
	// Skip leading slash
	for pos < len(path) && path[pos] == '/' {
		pos++
	}

	// Base case: consumed entire path
	if pos >= len(path) {
		if node.chainedHandler != nil {
			return node
		}
		// Check if there's a wildcard child that can match empty
		if node.wildcardChild != nil && node.wildcardChild.chainedHandler != nil {
			*params = append(*params, ParamCapture{
				Name:  node.wildcardChild.paramName,
				Value: "",
			})
			return node.wildcardChild
		}
		return nil
	}

	// Find end of current segment
	segEnd := pos
	for segEnd < len(path) && path[segEnd] != '/' {
		segEnd++
	}
	segment := path[pos:segEnd]

	// Priority: static > param > wildcard

	// 1. Try static children first (most common case)
	for _, child := range node.staticChildren {
		if child.prefix == segment {
			if result := t.lookupRecursive(child, path, segEnd, params); result != nil {
				return result
			}
		}
	}

	// 2. Try param child
	if node.paramChild != nil {
		*params = append(*params, ParamCapture{
			Name:  node.paramChild.paramName,
			Value: segment,
		})
		if result := t.lookupRecursive(node.paramChild, path, segEnd, params); result != nil {
			return result
		}
		// Backtrack: remove the param if match failed
		*params = (*params)[:len(*params)-1]
	}

	// 3. Try wildcard child (matches all remaining path)
	if node.wildcardChild != nil {
		// Capture remaining path (from current segment to end)
		remaining := path[pos:]
		// Remove trailing slash if present
		remaining = strings.TrimSuffix(remaining, "/")

		*params = append(*params, ParamCapture{
			Name:  node.wildcardChild.paramName,
			Value: remaining,
		})

		if node.wildcardChild.chainedHandler != nil {
			return node.wildcardChild
		}
		// Backtrack
		*params = (*params)[:len(*params)-1]
	}

	return nil
}

// PathExistsAnyMethod checks if a path exists in any method tree
// Used for 405 Method Not Allowed detection
func (t *RadixTree) PathExistsAnyMethod(path string) bool {
	for _, root := range t.trees {
		if root == nil {
			continue
		}

		// Check if path matches (ignoring whether handler exists)
		if t.pathExists(root, path, 0) {
			return true
		}
	}
	return false
}

// pathExists checks if a path exists in the tree (with or without handler)
func (t *RadixTree) pathExists(node *RadixNode, path string, pos int) bool {
	// Skip leading slashes
	for pos < len(path) && path[pos] == '/' {
		pos++
	}

	// Base case: consumed entire path
	if pos >= len(path) {
		// Path exists if this node has a handler (for any method)
		return node.handler != nil || node.chainedHandler != nil
	}

	// Find end of current segment
	segEnd := pos
	for segEnd < len(path) && path[segEnd] != '/' {
		segEnd++
	}
	segment := path[pos:segEnd]

	// Try static children
	for _, child := range node.staticChildren {
		if child.prefix == segment {
			if t.pathExists(child, path, segEnd) {
				return true
			}
		}
	}

	// Try param child
	if node.paramChild != nil {
		if t.pathExists(node.paramChild, path, segEnd) {
			return true
		}
	}

	// Try wildcard child
	if node.wildcardChild != nil {
		if node.wildcardChild.handler != nil || node.wildcardChild.chainedHandler != nil {
			return true
		}
	}

	return false
}
