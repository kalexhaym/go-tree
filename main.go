package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func walk(
	out io.Writer,
	path string,
	deep int,
	printFiles bool,
	parentLast bool,
	rowPrefix string,
) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	dirFiles, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	var files []os.DirEntry

	if !printFiles {
		for _, file := range dirFiles {
			if file.IsDir() {
				files = append(files, file)
			}
		}
	} else {
		files = dirFiles
	}

	filesLen := len(files) - 1

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	if deep > 0 {
		if !parentLast {
			rowPrefix += "│"
		}
		rowPrefix += "\t"
	}

	for index, file := range files {
		var err error

		isLast := index == filesLen

		filePrefix := "├───"
		if isLast {
			filePrefix = "└───"
		}

		if file.IsDir() {
			_, err = fmt.Fprintf(out, rowPrefix+filePrefix+file.Name()+"\n")
			if err != nil {
				return err
			}

			err = walk(
				out,
				filepath.Join(path, file.Name()),
				deep+1,
				printFiles,
				isLast,
				rowPrefix,
			)
			if err != nil {
				return err
			}
		} else if printFiles {
			fileInfo, err := file.Info()
			if err != nil {
				return err
			}

			fileName := fileInfo.Name()
			fileSize := fileInfo.Size()

			fileSizeText := " (empty)"
			if fileSize > 0 {
				fileSizeText = " (" + strconv.Itoa(int(fileSize)) + "b)"
			}

			_, err = fmt.Fprintf(out, rowPrefix+filePrefix+fileName+fileSizeText+"\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := walk(
		out,
		path,
		0,
		printFiles,
		false,
		"",
	)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
