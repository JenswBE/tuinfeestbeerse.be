#!/usr/bin/env bash

# Helper to rerender site when files changed.
# Execute with ./melthon_watch.sh [build options]

# Check if inotifywait is available
# https://stackoverflow.com/questions/592620/how-can-i-check-if-a-program-exists-from-a-bash-script
command -v inotifywait >/dev/null 2>&1 || { echo >&2 "inotifywait is missing, please install with 'sudo apt install inotify-tools'"; exit 1; }

# Start watching
echo "'melthon build ${*}' will be executed on file changes."
echo "Start watching ..."
inotifywait -qmr --exclude "output/.*" --event close_write --format "build ${*}" . | xargs -n1 melthon
