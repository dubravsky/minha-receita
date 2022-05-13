package check

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func checkZipFile(pth string) error {
	r, err := zip.OpenReader(pth)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", pth, err)
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return fmt.Errorf("error opening %s in %s: %w", f.Name, pth, err)
		}
		s := bufio.NewScanner(r)
		for s.Scan() {
			continue
		}
		if err := s.Err(); err != nil {
			return fmt.Errorf("error reading %s in %s: %w", f.Name, pth, err)
		}
		r.Close()
	}
	return nil
}

func checkZipFiles(dir string) (map[string]error, error) {
	r := make(map[string]error)
	ls, err := filepath.Glob(filepath.Join(dir, "*.zip"))
	if err != nil {
		return r, fmt.Errorf("error listing zip files: %w", err)
	}
	if len(ls) == 0 {
		return r, fmt.Errorf("no zip files found")
	}
	err = log.Output(2, fmt.Sprintf("Checking %d files…\n", len(ls)))
	if err != nil {
		return r, fmt.Errorf("error logging: %w", err)
	}
	var wg sync.WaitGroup
	for _, pth := range ls {
		wg.Add(1)
		go func(pth string) {
			defer wg.Done()
			err := checkZipFile(pth)
			if err != nil {
				log.Output(2, fmt.Sprintf("%s\tFAILED with\t%s", pth, err))
				r[pth] = err
			}
		}(pth)
	}
	wg.Wait()
	return r, nil
}

func Check(dir string, del bool) error {
	fails, err := checkZipFiles(dir)
	if err != nil {
		return fmt.Errorf("error checking zip files in %s: %w", dir, err)
	}
	if len(fails) != 0 {
		if del {
			for f := range fails {
				log.Output(2, fmt.Sprintf("Deleting %s", f))
				os.Remove(f)
			}
			return nil
		}
		return errors.New("error checking the zip files above")
	}
	return nil
}
