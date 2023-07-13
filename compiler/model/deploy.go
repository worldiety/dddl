package model

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func Deploy(dst string, src fs.FS) error {
	return fs.WalkDir(src, ".", func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		dir := path.Dir(p)
		if err := os.MkdirAll(filepath.Join(dst, dir), os.ModePerm); err != nil {
			return err
		}

		f, err := src.Open(p)
		if err != nil {
			return err
		}

		defer f.Close()

		buf, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(dst, p), buf, os.ModePerm); err != nil {
			return err
		}

		return nil

	})
}
