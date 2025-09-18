package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	dataFile string
	noColor  bool
)

var rootCmd = &cobra.Command {
	Use:   "memo",
	Short: "本地备忘录 CLI",
	Long:  `memo 是一个轻量级本地备忘录工具，所有数据保存在 JSON 文件。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	//全局flag
	home, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		return
	}
	defaultFile := filepath.Join(home, ".memo.json")
	//获取根命令rootCmd 的 PersistentFlags 集合
	rootCmd.PersistentFlags().StringVar(&dataFile, "data", defaultFile, "数据文件路径")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "关闭彩色输出")
}

// colorWrap 在 --no-color 时返回原字符串
func colorWrap(code, msg string) string {
	if noColor {
		return msg
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", code, msg)
}

// Success/Error 方便其他子命令复用
func Success(msg string) { fmt.Println(colorWrap("32", "✅ "+msg)) }
func Error(msg string)   { fmt.Println(colorWrap("31", "❌ "+msg)) }