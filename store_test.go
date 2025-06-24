package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	key         = "SoulSoceity"
	fileContent = "Zanka no Tachi"
	opts        = StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
)

func teadDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	expectedPath := "84647/b2184/badb9/331c9/ec324/5ec5f/aebce/3c140"
	expectedFilename := "84647b2184badb9331c9ec3245ec5faebce3c140"

	pathKey := CASPathTransformFunc(key)

	if pathKey.Path != expectedPath {
		t.Errorf("have %s want %s", pathKey.Path, expectedPath)
	}

	if pathKey.Filename != expectedFilename {
		t.Errorf("test 1have %s want %s", pathKey.Filename, expectedFilename)
	}
}

func TestStoreDelete(t *testing.T) {
	s := NewStore(opts)
	data := []byte(fileContent)

	if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	s := NewStore(opts)
	defer teadDown(t, s)

	// test write
	data := []byte(fileContent)
	if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}

	// test read
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, data, b)
}
