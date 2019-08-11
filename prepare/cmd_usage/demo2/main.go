package main

import (
	"fmt"
	"os/exec"
)

func main()  {
	var(
		cmd *exec.Cmd
		err error
		output []byte
	)

	cmd = exec.Command("/bin/bash","-c","ls -l")

	if output,err = cmd.CombinedOutput();err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}