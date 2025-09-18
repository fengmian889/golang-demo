package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"memo-demo/store"
	"os"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "删除所有已归档（done=true）的备忘录",
	Example: `memo clean`,
	Run: runClean,
}

func init() { rootCmd.AddCommand(cleanCmd) }

func runClean(cmd *cobra.Command, args []string) {
	entries, err := store.LoadEntries(dataFile)
	if err != nil {
		Error("读取数据失败: " + err.Error())
		os.Exit(1)
	}
	left := entries[:0]
	deleted := 0
	for _, e := range entries {
		if e.Done {
			deleted++
			continue
		}
		left = append(left, e)
	}
	if deleted == 0 {
		fmt.Println("没有已归档的备忘录需要清理。")
		return
	}

	if err := store.SaveEntries(dataFile, left); err != nil {
		Error("保存失败: " + err.Error())
		os.Exit(1)
	}
	Success(fmt.Sprintf("已清理 %d 条已归档备忘录", deleted))
}