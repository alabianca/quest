package trie

type TrieNode struct {
	Children  map[byte]*TrieNode
	EndOfWord bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children:  make(map[byte]*TrieNode),
		EndOfWord: false,
	}
}
