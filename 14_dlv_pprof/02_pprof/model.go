package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	From     string        `json:"from"`
	To       string        `json:"to"`
	Value    int64         `json:"value"`
	Metadata BlockMetadata `json:"metadata"`
}

type BlockMetadata struct {
	Nonce int64  `json:"nonce"`
	Hash  string `json:"hash"`
}

func (t Block) Hash() string {
	v := sha256.Sum256([]byte(fmt.Sprintf("%s->%s:%d:%d", t.From, t.To, t.Value, t.Metadata.Nonce)))
	return fmt.Sprintf("%x", v)
}
