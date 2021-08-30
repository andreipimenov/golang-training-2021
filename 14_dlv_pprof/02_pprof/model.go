package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
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
	buf := bytes.Buffer{}
	buf.WriteString(t.From)
	buf.WriteString("->")
	buf.WriteString(t.To)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(t.Value)))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(t.Metadata.Nonce)))

	v := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(v[:])
}
