package trie

import (
	"strings"
	"sync"
)

type Trie struct {
	Root  *TrieNode
	mutex *sync.RWMutex
}

func New() *Trie {
	return &Trie{
		Root:  NewTrieNode(),
		mutex: &sync.RWMutex{},
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

	// lock the trie for thread safe inserts
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, c := range lower {
		bt := byte(c)
		if current.Children[bt] == nil {
			current.Children[bt] = NewTrieNode()
		}
		current = current.Children[bt]
	}
	current.EndOfWord = true
}
