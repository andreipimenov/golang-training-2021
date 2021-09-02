package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

func Hash(from, to string, value, nonce int64) string {
	v := sha256.Sum256([]byte(fmt.Sprintf("%s->%s:%d:%d", from, to, value, nonce)))
	return fmt.Sprintf("%x", v)
}

func HashBuf(from, to string, value, nonce int64) string {
	buf := bytes.Buffer{}
	buf.WriteString(from)
	buf.WriteString("->")
	buf.WriteString(to)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(value)))
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(nonce)))
	v := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(v[:])
}
