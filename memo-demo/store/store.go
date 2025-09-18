package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)


type Entry struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Done    bool      `json:"done"`
	Created time.Time `json:"created"`
}

// LoadEntries 从指定路径加载所有备忘录
func LoadEntries(dataFile string) ([]Entry, error) {
	file, err := os.Open(dataFile)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	if err := json.NewDecoder(file).Decode(&entries); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return entries, nil
}

// SaveEntries 保存所有备忘录到指定路径
func SaveEntries(dataFile string, entries []Entry) error {
	tmp := dataFile + ".tmp"
	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(entries); err != nil {
		file.Close()
		os.Remove(tmp)
		return err
	}
	file.Close()
	return os.Rename(tmp, dataFile)
}
