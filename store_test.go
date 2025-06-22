package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var key = "SoulSoceity"

func TestPathTransformFunc(t *testing.T) {
	expectedPath := "84647/b2184/badb9/331c9/ec324/5ec5f/aebce/3c140"
	expectedOriginal := "84647b2184badb9331c9ec3245ec5faebce3c140"

	pathKey := CASPathTransformFunc(key)

	if pathKey.Path != expectedPath {
		t.Errorf("have %s want %s", pathKey.Path, expectedPath)
	}

	if pathKey.Filename != expectedOriginal {
		t.Errorf("test 1have %s want %s", pathKey.Filename, expectedOriginal)
	}
}

func TestStore(t *testing.T) {
	fileContent := "Zanka no Tachi"
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)

	// test write
	data := []byte(fileContent)
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
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
