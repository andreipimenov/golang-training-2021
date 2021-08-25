package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func MiningHandler(difficulty int, timeout time.Duration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		block := Block{}
		err = json.Unmarshal(body, &block)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		md := make(chan BlockMetadata)
		go mineBlock(block, difficulty, md)

		select {
		case v := <-md:
			block.Metadata.Nonce, block.Metadata.Hash = v.Nonce, v.Hash

			res, err := json.Marshal(block)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return

		case <-time.After(timeout):
			w.WriteHeader(http.StatusRequestTimeout)
			return
		}
	})
}

func mineBlock(block Block, difficulty int, md chan<- BlockMetadata) {
	prefix := strings.Repeat("0", difficulty)
	for i := int64(0); ; i++ {
		block.Metadata.Nonce = i
		hash := block.Hash()
		if strings.HasPrefix(hash, prefix) {
			block.Metadata.Hash = hash
			md <- block.Metadata
			break
		}
	}
}
