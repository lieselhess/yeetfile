package server

import (
	"fmt"
	"golang.org/x/crypto/nacl/secretbox"
	"os"
	"yeetfile/src/backblaze"
	"yeetfile/src/utils"
)

func TestDownload() {
	auth, err := backblaze.B2AuthorizeAccount(
		os.Getenv("B2_BUCKET_KEY_ID"),
		os.Getenv("B2_BUCKET_KEY"))

	if err != nil {
		panic(err)
	}

	id := ""
	length := 2665
	password := []byte("topsecret")

	salt, err := auth.B2PartialDownloadById(id, length-utils.KEY_SIZE, length)
	key, _, err := utils.DeriveKey(password, salt)
	if err != nil {
		return
	}

	// ---------------
	// TODO: Add password validation step before downloading from B2
	// ---------------

	out, err := os.OpenFile("out.enc", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	start := 0
	var output []byte
	for start < length-utils.KEY_SIZE-1 {
		chunkSize := utils.NONCE_SIZE + utils.BUFFER_SIZE + secretbox.Overhead + start
		if start+chunkSize > length-utils.KEY_SIZE-1 {
			chunkSize = length - utils.KEY_SIZE - 1
		}

		data, _ := auth.B2PartialDownloadById(id, start, chunkSize)

		plaintext, newStart := utils.DecryptChunk(key, data)
		output = append(output, plaintext...)
		start = newStart
	}

	_, _ = out.Write(output)
	_ = out.Close()

	plaintext, _ := os.ReadFile("out.enc")
	fmt.Println(string(plaintext))
}
