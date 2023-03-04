package hw10programoptimization

import (
	//"encoding/json"
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

// cd C:\REPO\Go\!OTUS\hwOTUS_YIA\hw10_program_optimization

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	emails, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(emails, domain)
}

// type users []User
type emails map[int64]string

//type emails []string

//type void struct{}

func getUsers(r io.Reader) (*emails, error) {
	/*
		content, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	*/
	result := make(emails)
	bufReader := bufio.NewReader(r)
	var i int64
	var email string
	var flag bool
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				//break
				flag = true
			} else {
				return nil, err
			}
		}
		email = gjson.Get(line, "Email").String()
		result[i] = email
		i++
		if flag {
			break
		}
	}
	//fmt.Println("result: ", result)
	//lines := strings.Split(string(content), "\n")
	//length := len(lines)
	//result := make(emails, length, length)
	/*
		var email string
		for i, line := range lines {
			email = gjson.Get(line, "Email").String()
			result[i] = email
		}
	*/
	return &result, nil
}

func countDomains(emails *emails, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var matched bool
	var err error
	var num int
	var domainstr string
	for _, email := range *emails {
		matched, err = regexp.Match("\\."+domain, []byte(email))
		if err != nil {
			return nil, err
		}

		if matched {
			domainstr = strings.ToLower(strings.SplitN(email, "@", 2)[1])
			num = result[domainstr]
			num++
			result[domainstr] = num
		}
	}
	return result, nil
}
