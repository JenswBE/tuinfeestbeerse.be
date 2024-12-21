#!/bin/bash

# Download test data

# Bash strict mode based on http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

# Settings
SCRIPT_DIR="$( dirname -- "$BASH_SOURCE"; )";

# Source helpers
. ${SCRIPT_DIR:?}/helpers.sh

# Push files
FOLDER_ID_ARTISTS=$(get_folder_id_by_name 'Artists')
FOLDER_ID_CAROUSEL=$(get_folder_id_by_name 'Carousel')
FOLDER_ID_SPONSORS=$(get_folder_id_by_name 'Sponsors')
push_file 'carousel-1.jpg' 'c530ce3b-cb68-48f0-9df1-3d77416bfe83' "${FOLDER_ID_CAROUSEL}"
push_file 'carousel-2.jpg' 'b719fec5-2f35-4768-a8af-016de6f03499' "${FOLDER_ID_CAROUSEL}"
push_file 'carousel-3.jpg' 'da3c604b-e0aa-4b6a-9ecc-eaecf161acf9' "${FOLDER_ID_CAROUSEL}"
push_file 'artist.png' '81f10ada-bacf-4056-85ee-520faaf34dcc' "${FOLDER_ID_ARTISTS}"
push_file 'sponsor-main.jpg' 'b38145ee-6cf3-4447-976f-225846269e3d' "${FOLDER_ID_SPONSORS}"
push_file 'sponsor-regular.png' 'a95192fc-a00b-49b5-8495-3144527f9093' "${FOLDER_ID_SPONSORS}"

# Push collections
push_collection "general_settings"
push_collection "carousel"
push_collection "carousel_files"
push_collection "sponsors"
push_collection "stages"
push_collection "artists"
push_collection "artist_links"
push_collection "performances"
push_collection "flow_debounce"

# Push users
push_local_user
