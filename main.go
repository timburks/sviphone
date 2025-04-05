package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

var transforms map[string]string = map[string]string{}

func run() error {
	err := filepath.WalkDir("events", func(path string, d os.DirEntry, err error) error {
		if strings.HasPrefix(path, "events/20") && !d.IsDir() {
			err = process(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = filepath.WalkDir("photos", func(path string, d os.DirEntry, err error) error {
		if strings.HasPrefix(path, "photos/20") && !d.IsDir() {
			err = process(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	//log.Printf("%+v", transforms)
	{
		b, err := os.ReadFile("events/events.html")
		if err != nil {
			return err
		}
		s := string(b)
		for k, v := range transforms {
			s = strings.ReplaceAll(s, k, v)
		}
		os.WriteFile("events.html", []byte(s), 0644)
	}
	{
		b, err := os.ReadFile("photos/photos.html")
		if err != nil {
			return err
		}
		s := string(b)
		for k, v := range transforms {
			s = strings.ReplaceAll(s, k, v)
		}
		os.WriteFile("photos.html", []byte(s), 0644)
	}
	return err
}

func process(path string) error {
	//log.Printf("%s", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	url := strings.TrimSpace(strings.TrimPrefix(lines[2], " url: "))
	//log.Printf("%s", url)
	transforms[url] = path
	return nil
}
