package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	templateDir := os.Getenv("HOME") + "/clipboard_templates"

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Printf("テンプレートディレクトリが見つかりません: %s\nディレクトリを作成してください。\n", templateDir)
		os.Exit(1)
	}

	currentClipboard, err := readClipboard()
	if err != nil {
		fmt.Println("クリップボードの内容を読み取れませんでした。")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		fmt.Println("テンプレートディレクトリの内容を読み取れませんでした。")
		os.Exit(1)
	}

	var templates []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			templates = append(templates, fmt.Sprintf("%s\t%s", file.Name(), filepath.Join(templateDir, file.Name())))
		}
	}

	if len(templates) == 0 {
		fmt.Printf("テンプレートが見つかりません。%s にテンプレートファイルを追加してください。\n", templateDir)
		os.Exit(1)
	}

	selected, err := selectTemplate(templates, currentClipboard)
	if err != nil {
		fmt.Println("テンプレートの選択中にエラーが発生しました。")
		os.Exit(1)
	}

	if selected == "" {
		fmt.Println("操作がキャンセルされました。")
		os.Exit(0)
	}

	selectedFile := strings.Split(selected, "\t")[1]

	appendString, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		fmt.Println("選択されたテンプレートファイルを読み取れませんでした。")
		os.Exit(1)
	}

	newClipboard := fmt.Sprintf("%s\n\n----\n%s\n----", appendString, currentClipboard)

	if err := writeClipboard(newClipboard); err != nil {
		fmt.Println("クリップボードの内容を更新できませんでした。")
		os.Exit(1)
	}

	fmt.Printf("\n【更新後のクリップボード内容】\n--------------------\n%s\n--------------------\n", newClipboard)
}

func selectTemplate(templates []string, currentClipboard string) (string, error) {
	cmd := exec.Command("fzf", "--delimiter", "\t", "--with-nth=1", "--preview", fmt.Sprintf("cat {2} && echo -e \"----\n%s\n----\n\" | head -70", currentClipboard), "--preview-window", "up:70%:wrap", "--prompt", "追加したいテンプレートを選択してください: ")
	cmd.Stdin = strings.NewReader(strings.Join(templates, "\n"))

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func readClipboard() (string, error) {
	cmd := exec.Command("pbpaste")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func writeClipboard(content string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(content)
	return cmd.Run()
}
