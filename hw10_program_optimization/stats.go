package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {

	result := make(DomainStat)
	bufReader := bufio.NewReader(r)
	var email string
	var flag bool
	var matched bool
	var err error
	var num int
	var domainstr string
	var line string
	searchDomain := "." + domain
	for {
		line, err = bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				flag = true
			} else {
				return nil, err
			}
		}
		email = gjson.Get(line, "Email").String()

		matched = strings.Contains(email, searchDomain)

		if matched {
			domainstr = strings.ToLower(strings.SplitN(email, "@", 2)[1])
			num = result[domainstr]
			num++
			result[domainstr] = num
		}
		if flag {
			break
		}
	}

	return result, nil
}
