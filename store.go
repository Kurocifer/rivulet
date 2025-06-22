package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CASPathTransformFunc takes a key descriptor for a file,
// and returns it's path on disk.
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5                       // determines the level of embeddin gof the path (the number of subdirectores)
	sliceLen := len(hashStr) / blockSize // So a hashStr len of 10, will yeild 2 sub directores
	path := make([]string, sliceLen)

	for i := range sliceLen {
		start, end := i*blockSize, (i*blockSize)+blockSize
		path[i] = hashStr[start:end]
	}

	fmt.Println("here's the hash string: " + hashStr)

	return PathKey{
		Path:     strings.Join(path, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Path     string
	Filename string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Path, p.Filename)
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc // Convert the key into the file path
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pahtKey := s.PathTransformFunc(key)
	return os.Open(pahtKey.FullPath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	// create the full path with all it's sub directories
	if err := os.MkdirAll(pathKey.Path, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := pathKey.FullPath()
	// Open desired file
	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read write connection payload into file
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s\n", n, pathAndFilename)
	return nil
}
