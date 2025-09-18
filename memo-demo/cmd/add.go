package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"memo-demo/store"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "新增一条备忘录",
	Example: `  memo add -t cobra -c "learn cobra"
  memo add --title demo --content "build a memo app" --data /tmp/m.json`,
	Run: runAdd,
}

var (
	addTitle   string
	addContent string
)

func init() {
	//添加addcmd子命令和参数
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addTitle, "title", "t", "", "备忘录标题（必填）")
	addCmd.Flags().StringVarP(&addContent, "content", "c", "", "备忘录内容（必填）")
	//将参数标记为必填项
	addCmd.MarkFlagRequired("title")
	addCmd.MarkFlagRequired("content")
}

func runAdd(cmd *cobra.Command, args []string){
	// 1. 读取现有数据
	entries,err := store.LoadEntries(dataFile)
	if err != nil {
		Error("读取数据失败: " + err.Error())
		os.Exit(1)
	}
	// 2. 标题冲突检查
	for _, e := range entries {
		if e.Title == addTitle {
			Error(fmt.Sprintf("标题 %q 已存在", addTitle))
			os.Exit(1)
		}
	}

	// 3. 追加新条目
	entries = append(entries, store.Entry{
		Title:   addTitle,
		Content: addContent,
		Done:    false,
		Created: time.Now(),
	})

	// 4. 写回文件
	if err := store.SaveEntries(dataFile,entries);err != nil {
		Error("保存失败: " + err.Error())
		os.Exit(1)
	}
	Success(fmt.Sprintf("已添加: %s", addTitle))
}
