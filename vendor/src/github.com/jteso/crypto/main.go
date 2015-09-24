// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// Crypto package contains some crypto utilities
package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// generateUUID simply returns a UUID
func GenerateUUID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %v", err)
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16]), nil
}

// computeMD5 returns the MD5 of the content
// of a file passed
func ComputeMd5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "Error", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "Error", err
	}

	checksum := hash.Sum(result)

	return hex.EncodeToString(checksum), nil
}
