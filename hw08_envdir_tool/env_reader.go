package main

import (
	"errors"
	"fmt"
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

var (
	ErrNotDirectory = errors.New("it's file, not directory!")
)

// cd  C:\REPO\Go\!OTUS\hwOTUS_YIA\hw08_envdir_tool

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
				fmt.Println("file: ", info.Name())
				if info.Size() > 0 {
					file, innererr := os.Open(path)
					if innererr != nil {
						return innererr
					}
					rez, innererr := readContent(file)
					if innererr != nil {
						return innererr
					}
					fmt.Println("rez: ", rez)
					envExempl.Value = clearFileContent(rez)
					innererr = file.Close()
					if innererr != nil {
						return innererr
					}
				} else {
					envExempl.NeedRemove = true
				}
				fmt.Println("exempl: ", envExempl)
				clearedName := clearFileName(info.Name())
				envir[clearedName] = envExempl
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	fmt.Println("func_out: ", envir)
	return envir, nil
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
		//fmt.Println("buf: ", buf)
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
	//output = strings.ReplaceAll(input, 0x00, "\n")
	return strings.TrimRight(input, " ")
}
