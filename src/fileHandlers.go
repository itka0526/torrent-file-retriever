package src

import (
	"fmt"
	"io/fs"
	"os"
)

func GetFiles() ([]fs.DirEntry, error) {
	p := "./downloads/"
	entries, err := os.ReadDir(p)

	fmt.Println(entries)
	if err != nil {
		return entries, err
	}

	return entries, nil
}
