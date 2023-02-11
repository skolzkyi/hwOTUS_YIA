package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	from, to             string
	limit, offset        int64
	ErrSourceIsNotExited = errors.New("source is not exited")
	ErrSourcePathIsNull  = errors.New("source path is null")
	ErrTargetPathIsNull  = errors.New("target path is null")
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	var err error
	sep := string(os.PathSeparator)

	err = checkInputFlags(from, to)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if from == to {
		toPathSlice := strings.Split(to, ".")
		toPathSlice[len(toPathSlice)-2] = toPathSlice[len(toPathSlice)-2] + "_copy"
		to = strings.Join(toPathSlice, ".")
	}
	toDirPathSlice := strings.Split(to, sep)
	toDirPathSlice = toDirPathSlice[:len(toDirPathSlice)-3]
	toDir := strings.Join(toDirPathSlice, sep)
	_, err = os.Stat(toDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(toDir, 0750)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		} else {
			fmt.Println(err.Error())
			return
		}
	}
	err = Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

//cd C:\REPO\Go\!OTUS\hwOTUS_YIA\hw07_file_copying

func checkInputFlags(from, to string) error {

	var err error
	fmt.Println("check")
	if from == "" {
		err = ErrSourcePathIsNull
		return err
	}

	if to == "" {
		err = ErrTargetPathIsNull
		fmt.Println(err.Error())
		return err
	}

	_, err = os.Stat(from)
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrSourceIsNotExited
			fmt.Println(err.Error())
			return err
		} else {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil

}
