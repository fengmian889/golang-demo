package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "对httpt服务器发起一个http请求",
	Long:  "对httpt服务器发起一个http请求，可以指定header，body等请求参数",
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		method, _ := cmd.Flags().GetString("method")
		headerStrings, _ := cmd.Flags().GetStringSlice("header")
		var headers http.Header
		if len(headerStrings) > 0 {
			headers := parseHeaders(headerStrings)
			fmt.Print(headers)
		}
		bodySource, _ := cmd.Flags().GetString("body")
		requestBody, err := getBodyContent(bodySource)
		if err != nil {
			log.Fatalf("处理请求体时出错: %v", err)
		}
		requestCall(url, method, headers, requestBody)
	},
}

func parseHeaders(headerStrings []string) http.Header {
	headers := make(http.Header)
	for _, header := range headerStrings {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers.Add(key, value)
		}
	}
	return headers
}

func getBodyContent(source string) ([]byte, error) {
	if source == "" {
		return nil, nil // 没有提供 --body，返回 nil
	}

	if strings.HasPrefix(source, "@") {
		filePath := source[1:] // 去掉开头的 '@'
		log.Printf("从文件读取请求体: %s", filePath)
		// 使用 os.ReadFile 读取整个文件内容
		return os.ReadFile(filePath)
	}

	// 3. 如果不是文件，则视为原始字符串
	log.Printf("使用原始字符串作为请求体")
	return []byte(source), nil
}

func requestCall(url string, method string, header http.Header, requestBody []byte) {
	queustUrl := url
	httpMethod := method
	headers := header
	body := requestBody
	req, err := http.NewRequest(method, queustUrl, bytes.NewReader(body))
	if err != nil {
		log.Fatalf("错误：无法创建请求: %v", err)
	}
	req.Header = headers
	client := &http.Client{Timeout: 10 * time.Second}
	log.Printf("发送请求: %s %s", httpMethod, queustUrl)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("错误：请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取并打印响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("错误：无法读取响应体: %v", err)
	}
	log.Printf("响应状态: %s", resp.Status)
	log.Printf("响应头: %v", resp.Header)

	const maxOutputSize = 1 << 20
	if len(respBody) > maxOutputSize {
		fmt.Printf("响应体: %s... (已截断，完整内容大小为 %d 字节)\n", string(respBody[:maxOutputSize]), len(respBody))
	} else {
		fmt.Printf("响应体: %s",string(respBody))
	}
}

func init() {
	rootCmd.AddCommand(callCmd)
	callCmd.Flags().StringP("method", "X", "GET", "指定服务器响应的 HTTP 方法")
	callCmd.Flags().StringSliceP("header", "H", []string{}, "可以多次添加的自定义请求头，格式为 'Key: Value'")
	callCmd.Flags().StringP("body", "d", "", "设置请求体，可以是字符串或文件路径（以@开头）")
}
