#!/bin/bash

SCRIPT=$(realpath "$0")
SCRIPTPATH=$(dirname "$SCRIPT")
pushd "$SCRIPTPATH/go" >/dev/null

# The directory to watch
WATCHED_DIR="./"

# The command to run when a change is detected
COMMAND="GOOS=js GOARCH=wasm go build -o ../static/dist.wasm ."

# Use inotifywait to watch the directory for any activity
# The -m flag means to monitor continuously
# The -e flag specifies the event to watch for; modify, create, delete are common events
# The --format option specifies the output format for events
# The %f outputs the file name, and %w the watched path

echo "Building..."
eval $COMMAND

inotifywait -r -m -e modify -e create -e delete --format '%w%f' "$WATCHED_DIR" | while read FILE
do
  # Run the command when a change is detected
  echo "Building..."
  eval $COMMAND
done