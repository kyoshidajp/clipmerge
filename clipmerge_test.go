package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	// Setup
	templateDir := filepath.Join(os.Getenv("HOME"), "clipboard_templates")
	os.MkdirAll(templateDir, os.ModePerm)
	defer os.RemoveAll(templateDir)

	templateFile := filepath.Join(templateDir, "template.txt")
	ioutil.WriteFile(templateFile, []byte("template content"), os.ModePerm)

	clipboard.WriteAll("current clipboard content")

	// Execute
	main()

	// Verify
	updatedClipboard, err := clipboard.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read clipboard content: %v", err)
	}
	expectedClipboard := "template content\n\n----\ncurrent clipboard content\n----"
	assert.Equal(t, expectedClipboard, updatedClipboard)
}

func TestReadClipboard(t *testing.T) {
	// Setup
	expectedContent := "clipboard content"
	clipboard.WriteAll(expectedContent)

	// Execute
	content, err := readClipboard()

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestGetTemplates(t *testing.T) {
	// Setup
	templateDir := filepath.Join(os.Getenv("HOME"), "clipboard_templates")
	os.MkdirAll(templateDir, os.ModePerm)
	defer os.RemoveAll(templateDir)

	templateFile := filepath.Join(templateDir, "template.txt")
	ioutil.WriteFile(templateFile, []byte("template content"), os.ModePerm)

	files, _ := ioutil.ReadDir(templateDir)

	// Execute
	templates := getTemplates(files, templateDir)

	// Verify
	expectedTemplates := []string{"template.txt\t" + templateFile}
	assert.Equal(t, expectedTemplates, templates)
}

func TestSelectTemplate(t *testing.T) {
	// Setup
	templates := []string{"template.txt\t/path/to/template.txt"}
	currentClipboard := "current clipboard content"

	cmd := exec.Command("echo", "template.txt\t/path/to/template.txt")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	// Execute
	selected, err := selectTemplate(templates, currentClipboard)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "template.txt\t/path/to/template.txt", selected)
}

func TestWriteClipboard(t *testing.T) {
	// Setup
	expectedContent := "new clipboard content"

	// Execute
	err := writeClipboard(expectedContent)

	// Verify
	assert.NoError(t, err)
	content, _ := clipboard.ReadAll()
	assert.Equal(t, expectedContent, content)
}

func TestLint(t *testing.T) {
	cmd := exec.Command("golint", "./...")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	assert.NoError(t, err, out.String())
}

func TestVet(t *testing.T) {
	cmd := exec.Command("go", "vet", "./...")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	assert.NoError(t, err, out.String())
}
