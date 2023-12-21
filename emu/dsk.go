// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package emu

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"retro/emu/config"
	"retro/emu/device/diskette"
	"retro/emu/virtual"
)

// insertDisks mounts disk images, if any, into drives in slot #6.
func insertDisks(conf *config.Config, bridge *virtual.Bridge) error {

	slot := bridge.Memory().Slot(6)
	card, ok := slot.(*diskette.Card)
	if !ok {
		return nil
	}

	paths := []string{
		conf.Disk.Drive1,
		conf.Disk.Drive2,
	}

	for i := 0; i < 2; i++ {
		if len(paths[i]) == 0 {
			continue
		}

		stream, err := openImageStream(paths[i], conf.Version)
		if err != nil {
			return err
		}

		drv := card.Drive(i)
		drv.Insert(diskette.NewStandardImage().MustLoad(stream))

		if closer, ok := stream.(io.Closer); ok {
			_ = closer.Close()
		}
	}
	return nil
}

func openImageStream(path string, userAgent string) (io.Reader, error) {
	stream, err := openRemoteImage(path, userAgent)
	if err != nil {
		return openLocalImage(path)
	}
	return stream, nil
}

func openLocalImage(path string) (io.Reader, error) {
	return os.Open(path)
}

func openRemoteImage(path string, userAgent string) (io.Reader, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("retro %s", userAgent))
	req.Header.Set("Accept-Encoding", "gzip")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Do not pass the body reader. Read all first.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(body)

	if res.Header.Get("Content-Encoding") == "gzip" {
		return gzip.NewReader(reader)
	}
	return reader, nil
}
