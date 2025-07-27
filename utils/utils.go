package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	crypt "github.com/kurocifer/rivulet/rivulet-base/crypto"
	"github.com/kurocifer/rivulet/rivulet-base/p2p"
	"github.com/kurocifer/rivulet/rivulet-base/server"
	"github.com/kurocifer/rivulet/rivulet-base/store"
)

const (
	EncryptionKey = "EncryptionKey"
	ID            = "ID"
)

// MakeServer, returns a new file server.
func MakeServer(listenAddr string, nodes ...string) *server.FileServer {
	var encryptionKey []byte
	var ID string

	stringKey, err := ReadFromFile(".keys", EncryptionKey)
	if err != nil {
		encryptionKey = crypt.NewEncryptionKey()
		WriteToFile(".keys", EncryptionKey, string(encryptionKey))
	} else {
		encryptionKey = []byte(stringKey)
	}

	ID, _ = ReadFromFile(".keys", ID)

	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.DefaultHandSake,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := server.FileServerOpts{
		ID:                ID,
		EncKey:            encryptionKey,
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: store.CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := server.NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer
	return s
}

// GetWorkDir, returns the working directory of the rivulet daemon.
func GetWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./rivulet/"
	}

	return homeDir + "/rivulet/"
}

// CreateWorkDir, creates teh rivulet working directory.
func CreateWorkDir() error {
	dir := GetWorkDir()
	return os.MkdirAll(dir+"/.daemon", 0755)
}

// WriteToFile, writes a key and value to filename in the rivulet working directory.
func WriteToFile(filename, key, value string) error {
	err := CreateWorkDir()
	if err != nil {
		return err
	}

	filePath := GetWorkDir() + filename
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = fmt.Fprintf(writer, "%s:%s\n", key, value)
	if err != nil {
		return fmt.Errorf("error while writing to %s: %w", filename, err)
	}

	return writer.Flush()
}

// ReadFromFile, reads and returns the value of key from filename.
func ReadFromFile(filename, key string) (string, error) {
	filePath := GetWorkDir() + filename
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to read line: %w", err)
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		fmt.Println(parts[0])

		if parts[0] == key {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("key '%s' not found in in file %s", key, filename)
}
