package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

// go run code_user/main.go
func main() {
	cmd := exec.Command("go", "run", "code_user/main.go")
	var out, stderr bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = &out
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	io.WriteString(stdinPipe, "1 2\n")
	//运行给定的code，比对拿到的输出结果和正确的输出结果
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err,stderr.String())
	}

	fmt.Println(out.String())
	fmt.Println(out.String() == "3\n")

}
