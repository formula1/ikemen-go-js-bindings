package main

import (
	"fmt"
	"ikemen-go-bindings/filesystem"
	"syscall/js"
)

var fs filesystem.AbstractFileSystem

func main() {
	fs = filesystem.NewFileSystem()

	fmt.Println("Hello, World! (from go)")
	js.Global().Get("console").Get("log").Invoke("Hello Invoke Console!")
	js.Global().Get("console").Call("log", "Hello Call Console!")
	JsMain()
	FsMain()

}
