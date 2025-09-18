package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memo-demo/store"
	"os"
)

var doneCmd = &cobra.Command{
	Use:     "done <title>",
	Short:   "把指定标题的条目标记为已完成（归档）",
	Example: `memo done cobra`,
	Args:    cobra.ExactArgs(1), // 必须且只能传 1 个 title
	Run:     runDone,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

func runDone(cmd *cobra.Command, args []string){
	title := args[0]
	entries,err := store.LoadEntries(dataFile)
	if err != nil {
		Error("读取数据失败: " + err.Error())
		os.Exit(1)
	}
	found := false
	for _,e := range entries {
		if e.Title == title {
			e.Done = true
			found = true
			break
		}
	}

	if !found {
		Error(fmt.Sprintf("未找到标题 %q", title))
		os.Exit(1)
	}

	if err := store.SaveEntries(dataFile, entries); err != nil {
		Error("保存失败: " + err.Error())
		os.Exit(1)
	}

	Success(fmt.Sprintf("已归档: %s", title))
}