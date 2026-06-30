package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fatal("Usage: gomodslice <entrypoint> <destination>")
	}

	entrypoint := os.Args[1]
	dest := os.Args[2]

	if _, err := os.Stat(entrypoint); err != nil {
		fatal("Entrypoint does not exist: %s", entrypoint)
	}

	entryAbs, err := filepath.Abs(entrypoint)
	check(err)

	destAbs, err := filepath.Abs(dest)
	check(err)

	if entryAbs == destAbs {
		fatal("Entrypoint and destination must be different")
	}

	module := strings.TrimSpace(run("go", "list", "-m", "-f", "{{.Path}}"))

	fmt.Println("Module:", module)

	check(os.MkdirAll(dest, 0755))

	copyFile("go.mod", filepath.Join(dest, "go.mod"))

	if _, err := os.Stat("go.sum"); err == nil {
		copyFile("go.sum", filepath.Join(dest, "go.sum"))
	}

	output := run(
		"go",
		"list",
		"-deps",
		"-f",
		`{{if and (not .Standard) (eq .Module.Path "`+module+`")}}{{.Dir}}|{{.ImportPath}}{{end}}`,
		entrypoint,
	)

	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		dir := parts[0]
		importPath := parts[1]

		rel := strings.TrimPrefix(importPath, module+"/")

		fmt.Println("Copying", rel)

		dst := filepath.Join(dest, rel)

		check(filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relative, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			target := filepath.Join(dst, relative)

			if info.IsDir() {
				return os.MkdirAll(target, 0755)
			}

			if filepath.Ext(path) != ".go" {
				return nil
			}

			copyFile(path, target)
			return nil
		}))
	}

	fmt.Println("Done.")
}

func copyFile(src, dst string) {
	check(os.MkdirAll(filepath.Dir(dst), 0755))

	in, err := os.Open(src)
	check(err)
	defer in.Close()

	out, err := os.Create(dst)
	check(err)
	defer out.Close()

	_, err = io.Copy(out, in)
	check(err)
}

func run(name string, args ...string) string {
	cmd := exec.Command(name, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	check(cmd.Run())

	return strings.TrimSpace(stdout.String())
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func check(err error) {
	if err != nil {
		fatal("%v", err)
	}
}
