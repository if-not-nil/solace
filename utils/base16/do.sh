#!/bin/bash
# requirements: flavours

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WHERE=$1

for theme in "$SCRIPT_DIR"/*.yaml; do
    theme_name=$(basename "$theme" .yaml)

    flavours build "$theme" "$SCRIPT_DIR/css.mustache" > "$WHERE/$theme_name.css"
    echo "$theme -> $WHERE/$theme_name.css"
done

