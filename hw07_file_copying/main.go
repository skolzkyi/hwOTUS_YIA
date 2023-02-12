package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
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

	err = Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println()
		fmt.Println("Copy process completed successfully!")
	}

}

//cd C:\REPO\Go\!OTUS\hwOTUS_YIA\hw07_file_copying
