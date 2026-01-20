package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var targets = []struct {
	os   string
	arch string
	ext  string
}{
	{"windows", "amd64", ".exe"},
	{"linux", "amd64", ""},
	{"darwin", "amd64", ""},
	{"darwin", "arm64", ""}, // Apple Silicon
}

func main() {
	distDir := "dist"
	if err := os.MkdirAll(distDir, 0755); err != nil {
		panic(err)
	}

	for _, t := range targets {
		filename := fmt.Sprintf("tr2rl-%s-%s%s", t.os, t.arch, t.ext)
		path := filepath.Join(distDir, filename)

		fmt.Printf("Building %s ...\n", filename)

		cmd := exec.Command("go", "build", "-o", path, ".")
		cmd.Env = append(os.Environ(),
			"GOOS="+t.os,
			"GOARCH="+t.arch,
		)

		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("Failed to build %s: %v\n%s\n", filename, err, out)
		}
	}
	fmt.Println("Done! Binaries are in 'dist/' directory.")
}
