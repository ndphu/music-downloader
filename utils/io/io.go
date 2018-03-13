package io

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
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
	for _, r := range []string{"?", ":", "/", "|", "\\", "*", "\"", "<", ">"} {
		output = strings.Replace(output, r, "-", -1)
	}
	return output
}

func ReadFromUrl(input *url.URL) ([]byte, error) {
	resp, err := http.Get(input.String())
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	log.Println("Reading body...")
	return ioutil.ReadAll(resp.Body)
}

func GetWithCookie(input *url.URL, cookie *http.Cookie) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", input.String(), nil)
	if err != nil {
		return []byte{}, err
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != 200 {
		return []byte{}, errors.New(fmt.Sprintf("Server response with invalid status code %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)

}

func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}

func GetAuthDir() string {
	return path.Join(GetHomeDir(), ".music-downloader", "auth")
}
