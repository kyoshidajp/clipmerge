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
	templateDir := filepath.Join(os.Getenv("HOME"), "clipboard_templates")

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Printf("Template directory not found: %s\nPlease create the directory.\n", templateDir)
		os.Exit(1)
	}

	currentClipboard, err := readClipboard()
	if err != nil {
		fmt.Println("Could not read clipboard content.")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(templateDir)
	if (err != nil) {
		fmt.Println("Could not read template directory content.")
		os.Exit(1)
	}

	templates := getTemplates(files, templateDir)

	if len(templates) == 0 {
		fmt.Printf("No templates found. %s Please add template files.\n", templateDir)
		os.Exit(1)
	}

	selected, err := selectTemplate(templates, currentClipboard)
	if err != nil {
		fmt.Println("Error occurred while selecting template.")
		os.Exit(1)
	}

	if selected == "" {
		fmt.Println("Operation cancelled.")
		os.Exit(0)
	}

	selectedFile := strings.Split(selected, "\t")[1]

	appendString, err := ioutil.ReadFile(selectedFile)
	if err != nil {
		fmt.Println("Could not read selected template file.")
		os.Exit(1)
	}

	newClipboard := fmt.Sprintf("%s\n\n----\n%s\n----", appendString, currentClipboard)

	if err := writeClipboard(newClipboard); err != nil {
		fmt.Println("Could not update clipboard content.")
		os.Exit(1)
	}

	fmt.Printf("\n[Updated clipboard content]\n--------------------\n%s\n--------------------\n", newClipboard)
}

func getTemplates(files []os.FileInfo, templateDir string) []string {
	var templates []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			templates = append(templates, fmt.Sprintf("%s\t%s", file.Name(), filepath.Join(templateDir, file.Name())))
		}
	}
	return templates
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
	return clipboard.ReadAll()
}

func writeClipboard(content string) error {
	return clipboard.WriteAll(content)
}
