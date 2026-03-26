package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// 1. 让你输入 commit 信息
	fmt.Print("请输入提交信息：")
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)

	if msg == "" {
		fmt.Println("❌ 提交信息不能为空！")
		return
	}

	// 2. 自动执行 git 命令
	runCmd("git", "add", ".")
	runCmd("git", "commit", "-m", msg)
	runCmd("git", "pull", "origin", "main", "--rebase")
	err := runCmd("git", "push", "origin", "main")

	if err == nil {
		fmt.Println("\n✅ 代码上传成功！")
	} else {
		fmt.Println("\n❌ 上传失败：", err)
	}
}

// 执行系统命令，并输出结果
func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
