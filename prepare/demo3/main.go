package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		resultChan chan *result
		err    error
		output []byte
		res *result
	)

	resultChan = make(chan *result,100)

	//定义一个chan为了输出结果

	//用于取消进程
	ctx, cancel = context.WithCancel(context.TODO())

	go func() {
		var (
			cmd    *exec.Cmd
		)

		//执行bash
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 5;ls -l")

		//执行输出
		output, err = cmd.CombinedOutput()
	}()

	time.Sleep(time.Second)

	resultChan <- &result{
		err:err,
		output:output,
	}

	//取消上下文
	cancel()

	//输出result
	res = <- resultChan

	fmt.Println(string(res.output))
}



