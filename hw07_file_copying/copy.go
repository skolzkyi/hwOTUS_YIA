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
	symbol        string
	procentScale  int64
	rawDataLength int64
	scale         int64
	pbLength      int
	textPosStart  int
	textLength    int
}

func (pb *ProgressBarAC05) Init(rawDataLength int64) {
	pb.pbLength = 20
	pb.textPosStart = 9
	pb.textLength = 4
	pb.symbol = "|"
	pb.rawDataLength = rawDataLength
	pb.scale = rawDataLength / int64(pb.pbLength)
	pb.procentScale = rawDataLength / int64(100)
}

func (pb *ProgressBarAC05) Rewrite(curDataSize int64) {
	// time.Sleep(time.Second) // for PB testing
	count := curDataSize / pb.scale
	procent := int(curDataSize / pb.procentScale)
	var sb strings.Builder
	sb.WriteString("\r")
	for i := 1; i < pb.pbLength; i++ {
		switch {
		case i < pb.textPosStart:
			if int64(i) <= count {
				sb.WriteString(pb.symbol)
			} else {
				sb.WriteString(" ")
			}
		case i == pb.textPosStart:
			sb.WriteString(strconv.Itoa(procent))
			sb.WriteString("%")
			i += pb.textLength
		case i > pb.textPosStart+pb.textLength:
			if int64(i) <= count {
				sb.WriteString(pb.symbol)
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

	err = checkTargetDirectory(toPath, sep)
	if err != nil {
		return err
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	fileFromStatSize, err := checkFileFromAndReturnSize(fileFrom, offset)
	if err != nil {
		return err
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	fileFrom.Seek(offset, io.SeekStart)
	var ProgBar ProgressBarAC05
	ProgBar.Init(fileFromStatSize)
	var flag bool
	for !flag {
		buffer := make([]byte, 512)
		n, errFF := fileFrom.Read(buffer)
		if n < len(buffer) {
			buffer = buffer[:n]
		}
		writedBytes += int64(len(buffer))
		if writedBytes >= limit && limit > 0 {
			buffer = buffer[:int64(len(buffer))-(writedBytes-limit)]
			flag = true
		}
		ProgBar.Rewrite(writedBytes)
		_, errFT := fileTo.Write(buffer)
		if errFT != nil {
			return errFT
		}
		if errFF != nil {
			if errors.Is(errFF, io.EOF) {
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
		}
		return err
	}

	return nil
}

func checkFileFromAndReturnSize(fileFrom *os.File, offset int64) (int64, error) {
	fileFromStat, err := fileFrom.Stat()
	if err != nil {
		return 0, err
	}

	if fileFromStat.Size() == 0 {
		return 0, ErrUnsupportedFile
	}

	if fileFromStat.Size() < offset {
		return 0, ErrOffsetExceedsFileSize
	}

	return fileFromStat.Size(), nil
}

func checkTargetDirectory(targetDirectoryPath string, sep string) error {
	toDirPathSlice := strings.Split(targetDirectoryPath, sep)
	var toDir string
	if len(toDirPathSlice) > 1 {
		toDirPathSlice = toDirPathSlice[:len(toDirPathSlice)-1]
		toDir = strings.Join(toDirPathSlice, sep)
	} else {
		return nil
	}
	_, err := os.Stat(toDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(toDir, 0o750)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
