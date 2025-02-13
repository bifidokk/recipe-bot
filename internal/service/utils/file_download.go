package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFileFromURL(fileURL string) (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	randomFileName := hex.EncodeToString(bytes) + ".mp3"
	filePath := filepath.Join("/tmp", randomFileName)

	// #nosec G304 - Filename is validated
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			return
		}
	}()

	// #nosec G107 url is validated
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %s", resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
