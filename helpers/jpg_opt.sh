#!/usr/bin/env bash

# Helper to optimize jpeg files
# Execute with ./jpg_opt.sh FILES

# Check if jpegoptim is available
command -v jpegoptim >/dev/null 2>&1 || { echo >&2 "jpegoptim is missing, please install with 'sudo apt install jpegoptim'"; exit 1; }

# Optimize images
jpegoptim -t -s --all-progressive -m75 $*
