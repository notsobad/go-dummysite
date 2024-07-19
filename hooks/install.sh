#!/bin/bash

pushd .git/hooks || { echo "No .git/hooks directory found"; exit 1; }

for file in ../../hooks/*; do
  ln -sf $file
done

echo "Hooks installed successfully"
