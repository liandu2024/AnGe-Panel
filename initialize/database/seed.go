package database

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"sun-panel/lib/cmn"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		// file
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return copyFile(p, target)
	})
}

// EnsureSeedUploads copies ./seed/uploads to ./uploads if uploads does not exist.
// This is safe (won't overwrite existing uploads).
func EnsureSeedUploads() {
	if ok, _ := cmn.PathExists("uploads"); ok {
		return
	}
	seed := filepath.Join("seed", "uploads")
	if ok, _ := cmn.PathExists(seed); !ok {
		return
	}
	if err := copyDir(seed, "uploads"); err != nil {
		log.Println("seed uploads copy failed:", err)
	}
}
