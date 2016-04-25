package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"io"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	var dlURL = flag.String("dl", "http://gd.tuwien.ac.at/linux/mint/isos//stable/17.3/linuxmint-17.3-cinnamon-64bit.iso", "download url")
	var upURL = flag.String("up", "", "upload url")

	if dlURL == nil {
		log.Fatalln("Don't do this to me now")
	}

	if upURL == nil {
		log.Fatalln("Don't do this to me now")
	}

	// RESP_START OMIT
	resp, err := http.Get(*dlURL)
	if err != nil {
		log.Fatalln(err)
	}
	// RESP_END OMIT

	// PIPE_START OMIT
	r, w := io.Pipe()

	eChan := make(chan error)

	go func() {
		_, err := http.Post(*upURL, "application/octet-stream", r)
		r.Close()
		eChan <- err
	}()
	// PIPE_END OMIT

	// ENCRY_START OMIT
	writer, err := encryptionWriter([]byte{0, 1, 2, 3}, w)
	if err != nil {
		log.Fatalln(err)
	}

	// ENCRY_END OMIT
	// FINISH_START OMIT
	_, err = io.Copy(writer, resp.Body)
	resp.Body.Close()
	writer.Close()
	if err != nil {
		log.Fatalln(err)
	}

	err = <-eChan
	log.Fatalln(err)
	// FINISH_END OMIT
}

func encryptionWriter(key []byte, w io.Writer) (io.WriteCloser, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, block.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	w.Write(iv)
	return cipher.StreamWriter{S: cipher.NewCTR(block, iv), W: w}, nil
}
