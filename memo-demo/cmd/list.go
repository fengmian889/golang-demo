package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memo-demo/store"
	"os"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "显示备忘录",
	Example: ` memo list -a -k test`,
	Run: runList,
}

var (
	keyword string
	all     bool
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&keyword, "keyword", "k", "", "关键词过滤（选填）")
	listCmd.Flags().BoolVarP(&all, "all", "a",false, "同时显示已归档")
}

func runList(cmd *cobra.Command, args []string) {
	kw := strings.ToLower(keyword)
	entries,err := store.LoadEntries(dataFile)
	if err != nil {
		Error("读取数据失败: " + err.Error())
		os.Exit(1)
	}

	shown := 0
	for _, e := range entries {
		//当显示归档flag为false且归档状态为false的时候，继续循环，不输出
		if !all && e.Done {
			continue
		}
		if kw != "" &&
			!strings.Contains(strings.ToLower(e.Title), kw) &&
			!strings.Contains(strings.ToLower(e.Content), kw) {
			continue
		}
		status := " "
		if e.Done {
			status = "✓"
		}
		fmt.Printf("[%s] %s: %s\n", status, e.Title, e.Content)
		shown++
	}

	if shown == 0 { // ← 友好提示
		fmt.Println("暂无备忘录。")
	}
}