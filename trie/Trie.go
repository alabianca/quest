package trie

import "strings"

type Trie struct {
	Root *TrieNode
}

func New() *Trie {
	return &Trie{
		Root: NewTrieNode(),
	}
}

func (t *Trie) Search(word string) bool {
	current := t.Root
	lower := strings.ToLower(word)
	for _, c := range lower {
		bt := byte(c)
		if current.Children[bt] == nil {
			return false
		}
		current = current.Children[bt]
	}
	return current.EndOfWord
}

func (t *Trie) StartsWith(word string) bool {
	current := t.Root
	lower := strings.ToLower(word)
	for _, c := range lower {
		bt := byte(c)
		if current.Children[bt] == nil {
			return false
		}
		current = current.Children[bt]
	}
	return true
}

func (t *Trie) Insert(word string) {
	current := t.Root
	lower := strings.ToLower(word)

	for _, c := range lower {
		bt := byte(c)
		if current.Children[bt] == nil {
			current.Children[bt] = NewTrieNode()
		}
		current = current.Children[bt]
	}
	current.EndOfWord = true
}
