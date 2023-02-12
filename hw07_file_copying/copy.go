package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSourceIsNotExisted    = errors.New("source is not existed")
	ErrSourcePathIsNull      = errors.New("source path is null")
	ErrTargetPathIsNull      = errors.New("target path is null")
	ErrBadLimit              = errors.New("limit less 0")
	ErrBadOffset             = errors.New("offset less 0")
)

type ProgressBarAC05 struct {
	procentScale  int64
	rawDataLength int64
	scale         int64
	symbol        string
	pbLength      int
	textPosStart  int
	curPBPos      int
	textLength    int
}

func (PB *ProgressBarAC05) Init(rawDataLength int64) {
	PB.pbLength = 20
	PB.textPosStart = 9
	PB.textLength = 4
	PB.symbol = "|"
	PB.rawDataLength = rawDataLength
	PB.scale = rawDataLength / int64(PB.pbLength)
	PB.procentScale = rawDataLength / int64(100)
}

func (PB *ProgressBarAC05) Rewrite(curDataSize int64) {
	//time.Sleep(time.Second) //for PB testing
	count := curDataSize / PB.scale
	procent := int(curDataSize / PB.procentScale)
	var sb strings.Builder
	sb.WriteString("\r")
	for i := 1; i < PB.pbLength; i++ {
		switch {
		case i < PB.textPosStart:
			if int64(i) <= count {
				sb.WriteString(PB.symbol)
			} else {
				sb.WriteString(" ")
			}
		case i == PB.textPosStart:
			sb.WriteString(strconv.Itoa(procent))
			sb.WriteString("%")
			i = i + PB.textLength
		case i > PB.textPosStart+PB.textLength:
			if int64(i) <= count {
				sb.WriteString(PB.symbol)
			} else {
				break
			}
		}
	}
	fmt.Print(sb.String())
}

func Copy(fromPath, toPath string, offset, limit int64) error {

	var writedBytes int64
	var err error

	sep := string(os.PathSeparator)

	err = checkInputFlags(fromPath, toPath, limit, offset)
	if err != nil {
		return err
	}

	if fromPath == toPath {
		toPathSlice := strings.Split(toPath, ".")
		toPathSlice[len(toPathSlice)-2] = toPathSlice[len(toPathSlice)-2] + "_copy"
		toPath = strings.Join(toPathSlice, ".")
	}
	toDirPathSlice := strings.Split(toPath, sep)
	//fmt.Println(toDirPathSlice)
	if len(toDirPathSlice) > 1 {
		toDirPathSlice = toDirPathSlice[:len(toDirPathSlice)-1]
		//fmt.Println(toDirPathSlice)
		toDir := strings.Join(toDirPathSlice, sep)
		//fmt.Println(toDir)
		//return nil
		_, err = os.Stat(toDir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(toDir, 0750)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	fileFromStat, err := fileFrom.Stat()
	if err != nil {
		return err
	}
	if fileFromStat.Size() == 0 {
		return ErrUnsupportedFile
	}
	if fileFromStat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	fileFrom.Seek(offset, io.SeekStart)
	var ProgBar ProgressBarAC05
	ProgBar.Init(fileFromStat.Size())
	var flag bool
	for !flag {
		buffer := make([]byte, 512, 512)
		n, errFF := fileFrom.Read(buffer)
		//fmt.Println(n)
		if n < len(buffer) {
			buffer = buffer[:n]
		}
		//n, errFF := io.ReadFull(fileFrom, buffer)
		writedBytes = writedBytes + int64(len(buffer))
		if writedBytes >= limit && limit > 0 {
			//fmt.Println("writedbytes, owerwrite: ", writedBytes, writedBytes-limit)
			buffer = buffer[:int64(len(buffer))-(writedBytes-limit)]
			//fmt.Println("new buffer length: ", len(buffer))
			flag = true
		}
		//fmt.Println("new buffer length2: ", len(buffer))
		ProgBar.Rewrite(writedBytes)
		_, errFT := fileTo.Write(buffer)
		//fmt.Println("writed: ", k)
		if errFT != nil {
			return errFT
		}
		if errFF != nil {
			if errors.Is(errFF, io.EOF) { //|| errors.Is(errFF, io.ErrUnexpectedEOF)      && n == 0
				//fmt.Println("exit: ", errFF.Error())
				flag = true
			} else {
				return errFF
			}
		}
	}

	return nil
}

func checkInputFlags(from, to string, limit, offset int64) error {

	var err error
	if from == "" {
		err = ErrSourcePathIsNull
		return err
	}

	if to == "" {
		err = ErrTargetPathIsNull
		return err
	}

	if limit < 0 {
		err = ErrBadLimit
		return err
	}

	if offset < 0 {
		err = ErrBadOffset
		return err
	}

	_, err = os.Stat(from)
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrSourceIsNotExisted
			return err
		} else {
			return err
		}
	}

	return nil

}
