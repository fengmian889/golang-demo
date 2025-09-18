package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memo-demo/store"
	"os"
)

var delCmd = &cobra.Command{
	Use:     "del <title>",
	Short:   "删除指定标题的备忘录",
	Example: `memo del cobra`,
	Args:    cobra.ExactArgs(1),
	Run:     runDel,
}

func init() { rootCmd.AddCommand(delCmd) }

func runDel(cmd *cobra.Command, args []string) {
	title := args[0]
	entries, err := store.LoadEntries(dataFile)
	if err != nil {
		Error("读取数据失败: " + err.Error())
		os.Exit(1)
	}
	//原地过滤切片，不用重新分配内存，但是达到清空底层数组的效果
	newList := entries[:0]
	found := false
	for _, e := range entries {
		if e.Title == title {
			found = true
			//要删除的就跳过
			continue
		}
		newList = append(newList, e)
	}
	if !found {
		Error(fmt.Sprintf("未找到标题 %q", title))
		os.Exit(1)
	}
	if err := store.SaveEntries(dataFile, newList); err != nil {
		Error("保存失败: " + err.Error())
		os.Exit(1)
	}
	Success(fmt.Sprintf("已删除: %s", title))
}
