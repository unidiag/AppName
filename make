#!/bin/bash

APPNAME="AppName"
APPNAME_LOWER="${APPNAME,,}"

if [ "$1" = "clean" ]; then
    echo "Cleaning build files..."

    rm -f "./$APPNAME_LOWER"
    rm -f "./$APPNAME_LOWER.db"
    rm -f ./*.sqlite3 2>/dev/null
    rm -f "./changelog.db"
    #rm -rf ./build/*

    echo "Done."
    exit 0
fi

sudo systemctl stop $APPNAME_LOWER.service

# get current version from main.go
CUR_VERSION=$(grep -oP 'const VERSION = "\K[0-9]+\.[0-9]+' main.go)

# increment version (0.01)
NEW_VERSION=$(LC_NUMERIC=C awk "BEGIN {printf \"%.2f\", $CUR_VERSION + 0.01}")

DATE=$(date +"%Y-%m-%d")
TIME=$(date +"%H:%M:%S")
UNIX=$(date +%s)

echo "Old version: $CUR_VERSION"
echo "New version: $NEW_VERSION"

echo
echo "Enter changes for version $NEW_VERSION (Ctrl+D to finish):"
TMP=$(mktemp)
nano $TMP
CHANGELOG=$(cat $TMP)
rm $TMP

echo "----------------------------------------" >> changelog.txt
echo "Version $NEW_VERSION - $DATE" >> changelog.txt
echo "$CHANGELOG" >> changelog.txt
echo >> changelog.txt

sed -i "s|const APPNAME = \".*\".*|const APPNAME = \"$APPNAME\"|g" main.go
sed -i "s|const VERSION = \".*\".*|const VERSION = \"$NEW_VERSION\"|g" main.go
sed -i "s|const BUILD_DATE = \".*\".*|const BUILD_DATE = \"$DATE\"|g" main.go
sed -i "s|const BUILD_TIME = \".*\".*|const BUILD_TIME = \"$TIME\"|g" main.go

cd ./frontend/

sed -i "s|REACT_APP_NAME *=.*|REACT_APP_NAME = $APPNAME|g" .env
sed -i "s|REACT_APP_VERSION *=.*|REACT_APP_VERSION = $NEW_VERSION|g" .env
sed -i "s|BUILD_UNIX *=.*|BUILD_UNIX = $UNIX|g" .env

npm run build

cd ../

echo '####################################'

go build -ldflags "-linkmode external -extldflags '-static'" -o $APPNAME_LOWER


if [ "$1" = "run" ]; then
    echo "Running binary..."
    ./$APPNAME_LOWER
    exit 0
fi

echo "Restart $APPNAME_LOWER.service..."
sudo systemctl restart $APPNAME_LOWER.service