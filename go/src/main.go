package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Hello, World! (from go)")
	js.Global().Get("console").Get("log").Invoke("Hello Invoke Console!")
	js.Global().Get("console").Call("log", "Hello Call Console!")
}
