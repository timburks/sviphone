package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
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
			err = removelinks(path)
			if err != nil {
				return err
			}
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
			err = removelinks(path)
			if err != nil {
				return err
			}
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
		removelinks("events.html")
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
		removelinks("photos.html")
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

func removelinks(filename string) error {
	var re1 = regexp.MustCompile(`href=http`)
	var re2 = regexp.MustCompile(`href="http`)
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	s := string(b)
	s = re1.ReplaceAllString(s, `ignore=`)
	s = re2.ReplaceAllString(s, `ignore="`)
	return os.WriteFile(filename, []byte(s), 0644)
}
