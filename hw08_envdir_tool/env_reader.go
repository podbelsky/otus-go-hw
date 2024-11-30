package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	newl  = 0x0a
	empty = 0x00
)

type Environment map[string]string

func (e Environment) Add(key, value string) {
	e[key] = value
}

func NewEnvironment() Environment {
	env := Environment(make(map[string]string))

	return env
}

func ReadDir(dir string) (Environment, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}

	env := NewEnvironment()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() || strings.ContainsRune(info.Name(), '=') {
			return nil
		}

		key := info.Name()
		value, e := readValue(path)
		if e != nil {
			return e
		}

		env.Add(key, value)

		return nil
	})

	return env, err
}

func readValue(path string) (string, error) {
	file, e := os.Open(path)
	defer file.Close() //nolint:staticcheck

	if e != nil {
		return "", e
	}

	r := bufio.NewReader(file)
	value, e := r.ReadBytes(newl)
	if e != nil && !errors.Is(e, io.EOF) {
		return "", e
	}

	value = bytes.ReplaceAll(value, []byte{empty}, []byte{newl})
	value = bytes.TrimRight(value, ` \t`)

	return string(value), nil
}
