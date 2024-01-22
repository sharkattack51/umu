package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/unicode/norm"
)

func main() {
	if len(os.Args) >= 2 {
		if IsExist(os.Args[1]) {
			var walkDirs []string
			_ = filepath.WalkDir(
				os.Args[1],
				func(path string, info fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					walkDirs = append(walkDirs, path)
					return nil
				})

			// reverse
			for i := 0; i < len(walkDirs)/2; i++ {
				walkDirs[i], walkDirs[len(walkDirs)-i-1] = walkDirs[len(walkDirs)-i-1], walkDirs[i]
			}

			for i := 0; i < len(walkDirs); i++ {
				path := walkDirs[i]
				dir, file := filepath.Split(path)

				// ignore dot file
				if file[0] == '.' {
					continue
				}

				// rename utf8mac -> utf8
				isNFC, strNorm := NormalizeNFD2NFC(file)
				if !isNFC {
					newPath := filepath.Join(dir, strNorm)
					_ = os.Rename(path, newPath)
					Println("utf8mac -> utf8", path+" -> "+newPath)
				}

			}
		}
	} else {
		Println("help", "Traversal and rename the specified root directory.")
	}
}

func NormalizeNFD2NFC(str string) (bool, string) {
	isNFC := norm.NFC.IsNormalString(str)
	strNorm := string(norm.NFC.Bytes([]byte(str)))
	return isNFC, strNorm
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Println(stat string, arg string) {
	fmt.Println(fmt.Sprintf("umu: [%v] %v", stat, arg))
}
