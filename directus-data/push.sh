#!/bin/bash

# Download test data

# Bash strict mode based on http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail

# Start a session
BASE_URL=http://localhost:8055
CURL_OPTS="--fail-with-body --silent --show-error"
SESSION_COOKIE=$(curl $CURL_OPTS -H "Content-Type: application/json" -c - -d '{"email":"admin@example.com", "password": "admin", "mode": "session"}' "$BASE_URL/auth/login")
SCRIPT_DIR="$( dirname -- "$BASH_SOURCE"; )";

# Push files
FOLDERS=$(curl $CURL_OPTS -b <(echo "$SESSION_COOKIE") "$BASE_URL/folders")
FOLDER_ID_ARTISTS=$(jq -r '.data[] | select(.name=="Artists") | .id' <(echo "$FOLDERS"))
FOLDER_ID_CAROUSEL=$(jq -r '.data[] | select(.name=="Carousel") | .id' <(echo "$FOLDERS"))
FOLDER_ID_SPONSORS=$(jq -r '.data[] | select(.name=="Sponsors") | .id' <(echo "$FOLDERS"))
push_file () {
  local FILENAME="$1"
  local FILE_PATH="${SCRIPT_DIR}/files/${FILENAME}"
  local FILE_ID="$2"
  local FOLDER_ID="$3"
  echo "Uploading file ${FILE_PATH} with ID ${FILE_ID} to folder ${FOLDER_ID} ..."
  curl $CURL_OPTS -X PATCH -b <(echo "$SESSION_COOKIE") -o /dev/null -F id=${FILE_ID} -F folder="${FOLDER_ID}" -F file="@${FILE_PATH}" "$BASE_URL/files/${FILE_ID}"
}
push_file 'carousel-1.jpg' 'c530ce3b-cb68-48f0-9df1-3d77416bfe83' "${FOLDER_ID_CAROUSEL}"
push_file 'carousel-2.jpg' 'b719fec5-2f35-4768-a8af-016de6f03499' "${FOLDER_ID_CAROUSEL}"
push_file 'carousel-3.jpg' 'da3c604b-e0aa-4b6a-9ecc-eaecf161acf9' "${FOLDER_ID_CAROUSEL}"
push_file 'artist.png' '81f10ada-bacf-4056-85ee-520faaf34dcc' "${FOLDER_ID_ARTISTS}"
push_file 'sponsor-main.jpg' 'b38145ee-6cf3-4447-976f-225846269e3d' "${FOLDER_ID_SPONSORS}"
push_file 'sponsor-regular.png' 'a95192fc-a00b-49b5-8495-3144527f9093' "${FOLDER_ID_SPONSORS}"

# Push collections
push_collection () {
  local COLLECTION_NAME="$1"
  local FILE_PATH="${SCRIPT_DIR}/${COLLECTION_NAME}.json"
  echo "Importing collection ${COLLECTION_NAME} from file ${FILE_PATH} ..."
  curl $CURL_OPTS -b <(echo "$SESSION_COOKIE") -F file="@${FILE_PATH};type=application/json" "$BASE_URL/utils/import/${COLLECTION_NAME}"
}
push_collection "general_settings"
push_collection "carousel"
push_collection "carousel_files"
push_collection "sponsors"
push_collection "stages"
push_collection "artists"
push_collection "artist_links"
push_collection "performances"
push_collection "flow_debounce"
