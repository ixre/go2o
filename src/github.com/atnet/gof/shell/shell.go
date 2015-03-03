package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	isDebug = false
)

func SetDebug(debug bool) {
	isDebug = debug
}

// Shell标准输出缓冲区
// 用于返回输出的内容
type shellStdBuffer struct {
	writer io.Writer
	buf    *bytes.Buffer
}

func newShellStdBuffer(writer io.Writer) *shellStdBuffer {
	return &shellStdBuffer{
		writer: writer,
		buf:    bytes.NewBuffer([]byte{}),
	}
}
func (this *shellStdBuffer) Write(p []byte) (n int, err error) {
	n, err = this.buf.Write(p)
	if this.writer != nil {
		n, err = this.writer.Write(p)
	}
	return n, err
}
func (this *shellStdBuffer) String() string {
	return string(this.buf.Bytes())
}

// 执行Shell命令
// 如果没有返回error,则命令执行成功，反之失败
// code返回命令执行返回的状态码,返回0表示执行成功
// output返回命令输出内容
func execCommand(command string, std_in io.Reader, std_out io.Writer,
	std_err io.Writer, debug bool) (code int, output string, err error) {

	var status syscall.WaitStatus //执行状态
	//var output string             //输出内容
	var stdout *shellStdBuffer //标准输出
	var stderr *shellStdBuffer //标准错误输出

	if strings.TrimSpace(command) == "" {
		return -1, "", errors.New("no such command")
	}

	if debug {
		fmt.Println(fmt.Sprintf("[COMMAND]:\n%s\n%s",
			command, strings.Repeat("-", len(command))))
	}

	var arr []string = strings.Split(command, " ")
	var cmd *exec.Cmd = exec.Command(arr[0], arr[1:]...)

	stdout = newShellStdBuffer(std_out)
	stderr = newShellStdBuffer(std_err)

	cmd.Stdout = stdout
	cmd.Stdin = std_in
	cmd.Stderr = stderr

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	status = cmd.ProcessState.Sys().(syscall.WaitStatus)
	isSuccess := cmd.ProcessState.Success()
	if debug {
		fmt.Println(strings.Repeat("-", len(command)))
		if isSuccess {
			fmt.Println("[OK] Status:", status.ExitStatus(),
				" Used Time:", cmd.ProcessState.UserTime(), "\n")
		} else {
			fmt.Println("[Fail] Status:", status.ExitStatus(),
				" Used Time:", cmd.ProcessState.UserTime(), "\n")
		}
	}

	if isSuccess {
		output = stdout.String()
	} else {
		output = stderr.String()
	}

	return status.ExitStatus(), output, nil
}

// 执行Shell命令
// 如果没有返回error,则命令执行成功，反之失败
// code返回命令执行返回的状态码,返回0表示执行成功
func Run(command string) (code int, output string, err error) {
	//return execCommand(command, os.Stdin, os.Stdout, os.Stdin, isDebug)
	return execCommand(command, nil, nil, nil, isDebug)
}

// 执行Shell命令，并输出到os.StdOut
// 错误输出到os.StdErr
func StdRun(command string) (code int, output string, err error) {
	return execCommand(command, os.Stdin, os.Stdout, os.Stdin, isDebug)
}

// (后台/静默)执行
// 仅仅执行命令，不需要捕获结果
func Brun(command string) (err error) {
	if strings.TrimSpace(command) == "" {
		return errors.New("no such command")
	}

	var arr []string = strings.Split(command, " ")
	var cmd *exec.Cmd = exec.Command(arr[0], arr[1:]...)
	err = cmd.Start()

	if isDebug {
		fmt.Print(fmt.Sprintf("[COMMAND]:\n%s\n%s",
			command, strings.Repeat("-", len(command))))
		if err != nil {
			fmt.Print("[Error]:", err.Error())
		}
		fmt.Print("\n")
	}
	return err
}
