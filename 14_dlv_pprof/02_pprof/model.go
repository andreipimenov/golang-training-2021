package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
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
	//Sprintf - ~300 rps
	//strings.Builder - up to 380 rps
	//bytes.Buffer - up to 402 rps
	var sb bytes.Buffer
	sb.WriteString(t.From)
	sb.WriteString("->")
	sb.WriteString(t.To)
	sb.WriteRune(':')
	sb.WriteString(strconv.Itoa(int(t.Value)))
	sb.WriteRune(':')
	sb.WriteString(strconv.Itoa(int(t.Metadata.Nonce)))
	//v := sha256.Sum256([]byte(sb.String()))
	v := sha256.Sum256(sb.Bytes())
	return fmt.Sprintf("%x", v)
}
