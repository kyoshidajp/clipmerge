#!/bin/bash

TEMPLATE_DIR="$HOME/clipboard_templates"

if [[ ! -d "$TEMPLATE_DIR" ]]; then
    echo "テンプレートディレクトリが見つかりません: $TEMPLATE_DIR"
    echo "ディレクトリを作成してください。"
    exit 1
fi

CURRENT_CLIPBOARD=$(pbpaste)

TEMPLATES=()
for file in "$TEMPLATE_DIR"/*.txt; do
    # ファイルが存在しない場合はスキップ
    [[ -e "$file" ]] || continue
    # ファイル名を取得（拡張子を含む）
    filename=$(basename "$file")
    # 配列に "filename<TAB>filepath" を追加
    TEMPLATES+=("$filename"$'\t'"$file")
done

if [[ ${#TEMPLATES[@]} -eq 0 ]]; then
    echo "テンプレートが見つかりません。$TEMPLATE_DIR にテンプレートファイルを追加してください。"
    exit 1
fi

export CURRENT_CLIPBOARD

SELECTED=$(printf "%s\n" "${TEMPLATES[@]}" | \
    fzf --delimiter=$'\t' \
        --with-nth=1 \
        --preview='cat {2} && echo -e "----\n$CURRENT_CLIPBOARD\n----\n" | head -70' \
        --preview-window=up:70%:wrap \
        --prompt="追加したいテンプレートを選択してください: ")

if [[ -z "$SELECTED" ]]; then
    echo "操作がキャンセルされました。"
    exit 0
fi

SELECTED_FILE=$(echo "$SELECTED" | cut -f2)

APPEND_STRING=$(cat "$SELECTED_FILE")

NEW_CLIPBOARD=$(printf "%s\n\n----\n%s\n----" "$APPEND_STRING" "$CURRENT_CLIPBOARD")

echo "$NEW_CLIPBOARD" | pbcopy

echo -e "\n【更新後のクリップボード内容】\n--------------------"
echo "$NEW_CLIPBOARD"
echo -e "--------------------\n"
