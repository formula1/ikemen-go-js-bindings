package api

import (
	"fmt"
	"syscall/js"
)

func noArgPrint() {
	fmt.Println("Hello, World! (from common)")
}

func goPrint(value string) {
	fmt.Println("Hello go from common: ", value)
}

func jsPrint(this js.Value, args []js.Value) {
	fmt.Println("Hello js from common: ", args[0])
}
