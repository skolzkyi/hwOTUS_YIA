package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrNotDirectory = errors.New("it's file, not directory")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envir := make(Environment)
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		err = ErrNotDirectory
		return nil, err
	}
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				var envExempl EnvValue
				if info.Size() > 0 {
					rez, innererr := readFirstStringFromFile(path)
					if innererr != nil {
						return innererr
					}
					envExempl.Value = clearFileContent(rez)
				} else {
					envExempl.NeedRemove = true
				}
				clearedName := clearFileName(info.Name())
				envir[clearedName] = envExempl
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return envir, nil
}

func readFirstStringFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	rez, err := readContent(file)
	if err != nil {
		return "", err
	}

	return rez, nil
}

func readContent(file *os.File) (string, error) {
	rez := make([]byte, 0)
	for {
		buf := make([]byte, 1)
		_, innererr := file.Read(buf)
		if innererr != nil {
			if errors.Is(innererr, io.EOF) {
				break
			}
			return "", innererr
		}
		if buf[0] == 10 { // \n
			break
		}
		if buf[0] != 13 {
			if buf[0] != 0 {
				rez = append(rez, buf...)
			} else {
				rez = append(rez, 10)
			}
		}
	}
	return string(rez), nil
}

func clearFileName(input string) string {
	return strings.ReplaceAll(input, "=", "")
}

func clearFileContent(input string) string {
	return strings.TrimRight(input, " \t")
}
