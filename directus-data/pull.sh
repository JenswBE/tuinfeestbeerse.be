#!/bin/bash

# Download test data

# Bash strict mode based on http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

# Settings
SCRIPT_DIR="$( dirname -- "$BASH_SOURCE"; )";

# Source helpers
. ${SCRIPT_DIR:?}/helpers.sh

# Pull collections
pull_collection "carousel_files"
pull_collection "sponsors"
pull_collection "stages" "id,name"
pull_collection "artists" "id,name,artist_type,genre,description,picture,sort"
pull_collection "artist_links"
pull_collection "performances"
pull_collection "flow_debounce"

# Pull singletons
pull_singleton "general_settings"
pull_singleton "carousel" "id"
