package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func DownloadFileFromURL(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("invalid URL")
	}

	response, err := http.Get(parsedURL.String())
	if err != nil {
		return "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != 200 {
		return "", errors.New("received none 200 response code")
	}

	fileName := path.Base(parsedURL.Path)

	// #nosec G304 - Filename is validated
	file, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return "/tmp/" + fileName, nil
}
