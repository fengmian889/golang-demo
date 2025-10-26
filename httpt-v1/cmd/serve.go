package cmd

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动一个本地回显服务器",
	Long:  "启动一个本地 HTTP 服务器，它会将收到的所有请求信息（方法、路径、头、体）以 JSON 格式原样返回",
	Run: func(cmd *cobra.Command, args []string) {
		ip := "127.0.0.1:"
		port, _ := cmd.Flags().GetString("port")
		addr := ip + port

		server := &http.Server{
			Addr:    addr,
			Handler: http.HandlerFunc(echoHandler),
		}

		log.Printf("Starting server on %s", addr)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed to start: %v", err)
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	},
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	payload := map[string]interface{}{
		"method": r.Method,
		"url":    r.URL.Path,
		"header": r.Header,
		"body":   body,
	}

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "指定服务器监听的端口")
}
