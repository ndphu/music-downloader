package io

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func DownloadFileWithRetry(filepath string, fileUrl string, retry int) (err error) {
	try := 0

	for {
		try++
		err = DownloadFile(filepath, fileUrl)
		if err == nil {
			return err
		}
		if try == retry {
			return err
		} else {
			log.Printf("%v\n", err)
			log.Printf("Retrying... %d\n", try)
		}
	}
	return err
}

func DownloadFile(filepath string, fileUrl string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func CleanupFileName(input string) string {
	output := input
	for _, r := range []string{"?", ":"} {
		output = strings.Replace(output, r, "", -1)
	}
	return output
}
