package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Get will get the value associated with an env key.
func Get(key string) string {
	return os.Getenv(key)
}

func init() {
	var (
		err  error
		cwd  string
		pair []string
	)

	cwd, err = os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	env := openenv(cwd)

	for _, keyval := range strings.Split(env, "\n") {
		pair = strings.Split(keyval, "=")
		if err = os.Setenv(pair[0], pair[1]); err != nil {
			fmt.Println(err)
		}
	}
}

func openenv(dir string) string {
	file := filepath.Join(dir, ".env")
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(b)
}
