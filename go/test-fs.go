package main

import (
	"fmt"
	"strconv"
)

func FsMain() {
	fmt.Println("FS Main Start ==Hello World!==")

	files, error := fs.ReadDir("/")
	if error != nil {
		fmt.Println("Error reading root directory")
		return
	}

	fileLen := len(files)
	for i := 0; i < fileLen; i++ {
		fmt.Println("File " + strconv.Itoa(i) + " -> " + files[i])
	}

	fmt.Println("FS Main End")
}
