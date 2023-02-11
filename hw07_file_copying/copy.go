package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	var writedBytes int64

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	fmt.Println("check filefrom")
	defer fileFrom.Close()

	fileTo, err := os.Create(toPath)
	fmt.Println("check fileto1")
	if err != nil {
		return err
	}
	fmt.Println("check fileto2")
	defer fileTo.Close()

	buffer := make([]byte, 1024)
	fileFrom.Seek(offset, io.SeekStart)

	var flag bool
	for {
		n, errFF := fileFrom.Read(buffer)
		fmt.Println(n)
		if n < len(buffer) {
			buffer = buffer[:n]
		}
		//n, errFF := io.ReadFull(fileFrom, buffer)
		writedBytes = writedBytes + int64(n)
		if writedBytes >= limit && limit > 0 {
			fmt.Println("writedbytes, owerwrite: ", writedBytes, writedBytes-limit)
			buffer = buffer[:int64(len(buffer))-(writedBytes-limit)]
			fmt.Println("new buffer length: ", len(buffer))
			flag = true
		}
		_, errFT := fileTo.Write(buffer)
		if errFT != nil {
			return errFT
		}
		if errFF != nil {
			if errors.Is(errFF, io.EOF) { //|| errors.Is(errFF, io.ErrUnexpectedEOF)      && n == 0
				fmt.Println("exit: ", errFF.Error())
				flag = true
			} else {
				return errFF
			}
		}
		if flag {
			break
		}
	}

	return nil
}
