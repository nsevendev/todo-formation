package projectroot

import (
	"log"
	"os"
	"path/filepath"
)

var root string

func init() {
	root = find()
}

func find() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	log.Fatal("[FATAL] Impossible de trouver la racine du projet (go.mod manquant)")
	return ""
}

func Path(subpaths ...string) string {
	return filepath.Join(append([]string{root}, subpaths...)...)
}

func Root() string {
	return root
}
