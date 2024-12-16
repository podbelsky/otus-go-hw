package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`\A\w+\z`)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if ok := pattern.MatchString(domain); !ok {
		return nil, errors.New("bad domain name")
	}

	result := make(DomainStat)
	suffix := "." + domain

	user := &User{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), user); err != nil {
			return nil, err
		}

		email := user.Email

		if strings.HasSuffix(email, suffix) {
			if idx := strings.IndexRune(email, '@'); idx >= 0 {
				emailDomain := strings.ToLower(email[idx+1:])
				result[emailDomain]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return result, nil
}
