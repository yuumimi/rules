package trie

import (
	"errors"
	"strings"
)

var (
	ErrInvalidDomain = errors.New("invalid domain")
)

type Node struct {
	children map[string]*Node
	matched  bool
}

type Trie struct {
	root *Node
}

func splitDomain(domain string) ([]string, error) {
	if domain != "" && domain[len(domain)-1] == '.' {
		return nil, ErrInvalidDomain
	}

	parts := strings.Split(domain, ".")
	if len(parts) == 1 {
		if parts[0] == "" {
			return nil, ErrInvalidDomain
		}

		return parts, nil
	}

	for _, part := range parts[1:] {
		if part == "" {
			return nil, ErrInvalidDomain
		}
	}

	return parts, nil
}

func (t *Trie) Insert(domain string, full bool) error {
	parts, err := splitDomain(domain)
	if err != nil {
		return err
	}

	node := t.root

	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]

		if node.children == nil {
			return nil
		}

		n := node.children[part]

		if n == nil {
			n = &Node{
				children: map[string]*Node{},
			}

			node.children[part] = n
		}

		node = n
	}

	if full {
		node.matched = true
	} else {
		node.children = nil
	}

	return nil
}

func (t *Trie) Dump() []string {
	list := make([]string, 0, 1024)

	t.root.dump(&list, []string{})

	return list
}

func New() *Trie {
	return &Trie{
		root: &Node{
			children: map[string]*Node{},
			matched:  false,
		},
	}
}

func (t *Node) dump(list *[]string, parts []string) {
	if t.children == nil {
// 	    *list = append(*list, joinDomain(parts))                                // List
// 		*list = append(*list, joinDomain(append(parts, "+")))                   // Clash
// 		*list = append(*list, "DOMAIN-SUFFIX," + joinDomain(parts))             // Loon
// 		*list = append(*list, "host-suffix, " + joinDomain(parts) + ", proxy")  // Quantumult X
		*list = append(*list, joinDomain(append(parts, "")))                    // Surge

		return
	}

	if t.matched {
// 		*list = append(*list, joinDomain(parts))                                // List
// 		*list = append(*list, joinDomain(parts))                                // Clash
// 		*list = append(*list, "DOMAIN," + joinDomain(parts))                    // Loon
// 		*list = append(*list, "host, " + joinDomain(parts) + ", proxy")         // Quantumult X
		*list = append(*list, joinDomain(parts))                                // Surge
	}

	for k, v := range t.children {
		v.dump(list, append(parts, k))
	}
}

func joinDomain(parts []string) string {
	domain := ""
	index := len(parts) - 1

	for index > 0 {
		domain += parts[index]
		domain += "."

		index--
	}

	domain += parts[index]

	return domain
}
