package shell

import (
	"fmt"
	"testing"
)

func Test_bash(t *testing.T) {
	//SetDebug(true)
	handleOutput(Run("mkdir /home/testdir"))
	handleOutput(Run("touch /home/testdir/1"))
	handleOutput(Run("touch /home/testdir/2"))
	handleOutput(Run("ls /home/testdir"))
	handleOutput(Run("rm -rf /home/testdir"))
}

func handleOutput(code int, output string, err error) {
	fmt.Println("[Code]:", code, "\n", output)
	if err != nil {
		fmt.Println(err.Error())
	}
}
