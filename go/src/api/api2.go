package api

import (
	"fmt"
	"syscall/js"
)

func noArgPrint2() {
	fmt.Println("Hello, World! (from common)")
}

func goPrint2(value string) {
	fmt.Println("Hello go from common: ", value)
}

func jsPrint2(this js.Value, args []js.Value) {
	fmt.Println("Hello js from common: ", args[0])
}
