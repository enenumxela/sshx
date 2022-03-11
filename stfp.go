package sshx

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func (client *Client) MkdirAll(DIR string) (err error) {
	path := string(filepath.Separator)
	directories := strings.Split(DIR, path)

	for _, directory := range directories {
		path = filepath.Join(path, directory)

		_, err := client.SFTP.Lstat(path)
		if err == nil {
			continue
		}

		if err := client.SFTP.Mkdir(path); err != nil {
			break
		}
	}

	return
}

func (client *Client) Upload(SRC, DEST string) (err error) {
	directory, filename := filepath.Split(DEST)

	if err := client.SFTP.MkdirAll(directory); err != nil {
		return err
	}

	if filename != "" {
		SRCFile, err := os.Open(SRC)
		if err != nil {
			return err
		}

		defer SRCFile.Close()

		DESTFile, err := client.SFTP.OpenFile(DEST, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
		if err != nil {
			return err
		}

		defer DESTFile.Close()

		if _, err = io.Copy(DESTFile, SRCFile); err != nil {
			return err
		}
	}

	return
}

func (client *Client) Download(SRC, DEST string) error {
	if _, err := os.Stat(DEST); os.IsNotExist(err) {
		if err = os.MkdirAll(DEST, os.ModePerm); err != nil {
			return err
		}
	}

	files, err := client.SFTP.ReadDir(SRC)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)

		go func(file os.FileInfo) error {
			defer wg.Done()

			NewSRC := filepath.Join(SRC, file.Name())
			NewDEST := filepath.Join(DEST, file.Name())

			if file.IsDir() {
				if err := client.Download(NewSRC, NewDEST); err != nil {
					return err
				}
			} else {
				if err := client.DownloadFile(NewSRC, NewDEST); err != nil {
					return err
				}
			}

			return nil
		}(file)
	}

	wg.Wait()

	return nil
}

func (client *Client) DownloadFile(SRC, DEST string) (err error) {
	SRCFile, err := client.SFTP.OpenFile(SRC, (os.O_RDONLY))
	if err != nil {
		return
	}

	defer SRCFile.Close()

	DESTFile, err := os.Create(DEST)
	if err != nil {
		return
	}

	defer DESTFile.Close()

	if _, err = io.Copy(DESTFile, SRCFile); err != nil {
		return
	}

	return nil
}
