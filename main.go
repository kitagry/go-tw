package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type exitCode int

const (
	exitCodeOK  exitCode = 0
	exitCodeErr exitCode = 0
)

var efmRegexp = regexp.MustCompile(`^(\s+)([\w_]+\.go):(\d+):`)

func main() {
	code := run(context.Background())
	os.Exit(int(code))
}

func run(ctx context.Context) exitCode {
	gomodPath, err := searchFile("go.mod")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return exitCodeErr
	}
	rootPath, _ := path.Split(gomodPath)

	packageName, err := getPackageName(rootPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return exitCodeErr
	}

	args := append([]string{"test"}, os.Args[1:]...)
	cmd := exec.Command("go", args...)
	var buf bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	cmd.Run()

	sc := bufio.NewScanner(&buf)
	sc.Split(bufio.ScanLines)
	w := bufio.NewWriter(os.Stdout)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "--- FAIL: ") {
			_, err := w.WriteString(line + "\n")
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return exitCodeErr
			}
			continue
		}

		lines := []string{line}
		var packagePath string
		for sc.Scan() {
			line := sc.Text()
			if !strings.HasPrefix(line, "FAIL") || line == "FAIL" {
				lines = append(lines, line)
				continue
			}
			packageLine := strings.Fields(line)[1]
			packagePath = strings.Replace(packageLine, packageName, ".", 1)
			lines = append(lines, line)
			break
		}

		packagePath = path.Join(rootPath, packagePath)

		for _, line := range lines {
			newLine := efmRegexp.ReplaceAllString(line, fmt.Sprintf("$1%s/$2:$3:", packagePath))
			_, err := w.WriteString(newLine + "\n")
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return exitCodeErr
			}
		}
	}
	err = w.Flush()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return exitCodeErr
	}
	return exitCodeOK
}

func getPackageName(rootPath string) (string, error) {
	data, err := os.ReadFile(path.Join(rootPath, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	for _, line := range bytes.Split(data, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("module ")) {
			return string(line[len("module "):]), nil
		}
	}
	return "", fmt.Errorf("module name not found")
}

func searchFile(filename string) (path string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return _searchFile(filename, pwd)
}

func _searchFile(filename, pwd string) (string, error) {
	if pwd == "/" {
		return "", fmt.Errorf("%s not found", filename)
	}
	files, err := os.ReadDir(pwd)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.Name() == filename {
			return path.Join(pwd, filename), nil
		}
	}

	parentPath, _ := path.Split(pwd)
	return _searchFile(filename, parentPath)
}
