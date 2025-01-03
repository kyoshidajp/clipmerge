# clipmerge

`clipmerge` is a tool to merge clipboard content with templates. It allows you to select a template and append it to the current clipboard content.

## Installation

1. Clone the repository:

```sh
git clone https://github.com/kyoshidajp/clipmerge.git
cd clipmerge
```

2. Install the dependencies:

```sh
go mod tidy
```

3. Build the Go program:

```sh
go build -o clipmerge
```

## Usage

1. Create a directory for your templates:

```sh
mkdir -p ~/clipboard_templates
```

2. Add your template files (with `.txt` extension) to the `~/clipboard_templates` directory.

3. Run the `clipmerge` program:

```sh
./clipmerge
```

4. Follow the prompts to select a template and merge it with the current clipboard content.

## Notes

- The `clipmerge` program uses the `fzf` command for template selection. Make sure you have `fzf` installed on your system.
- The program reads and writes clipboard content using the `github.com/atotto/clipboard` package.
- The program is cross-platform and handles directory separators using the `filepath.Join` function.
