package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootDirectory = "rivulet"

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Path     string
	Filename string
}

func (p PathKey) FullPath(rootDir string) string {
	return fmt.Sprintf("%s/%s/%s", rootDir, p.Path, p.Filename)
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Path:     key,
		Filename: key,
	}
}

type StoreOpts struct {
	// Root is the root directory of the system
	Root              string
	PathTransformFunc PathTransformFunc // Convert the key into the file path
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootDirectory
	}
	return &Store{
		StoreOpts: opts,
	}
}

// CASPathTransformFunc takes a key descriptor for a file,
// and returns it's path on disk. (I should probably come up with a better name)
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5                       // determines the level of embedding of the path (the number of subdirectores)
	sliceLen := len(hashStr) / blockSize // So a hashStr len of 10, will yeild 2 sub directores
	path := make([]string, sliceLen)

	for i := range sliceLen {
		start, end := i*blockSize, (i*blockSize)+blockSize
		path[i] = hashStr[start:end]
	}

	return PathKey{
		Path:     strings.Join(path, "/"),
		Filename: hashStr,
	}
}

func (s *Store) Has(key string) bool {
	PathKey := s.PathTransformFunc(key)

	_, err := os.Stat(PathKey.FullPath(s.Root))
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	topDir := strings.Split(pathKey.FullPath(s.Root), "/")[:1]
	if len(topDir) == 0 {
		return fmt.Errorf("top directory doesn't exist (I actually don't know what to say in this error)")
	}
	if err := os.RemoveAll(strings.Join(topDir, "/")); err != nil {
		return err
	}

	fmt.Printf("deleted [%s] from disk\n", pathKey.Filename)
	return nil
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
	return os.Open(pahtKey.FullPath(s.Root))
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	// create the full path with all it's sub directories
	if err := os.MkdirAll(s.Root+"/"+pathKey.Path, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := pathKey.FullPath(s.Root)
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
