#!/bin/bash
set -eo pipefail

sigint_handler()
{
  kill $PID
  exit
}

trap sigint_handler SIGINT

while true; do
  echo "Building app"
  if go build -o eventsourcing-hack ./cmd/eventsourcing-hack/main.go; then
    echo "App built, starting server."

    ./eventsourcing-hack &

    # Remember process id so we can kill it when a file changes
    PID=$!

  else
    echo "Build failed."
    unset PID
  fi

  printf "Waiting for file changes...\n\n"


  # Hang until a file changes
  inotifywait -e modify -e create -e delete --exclude \.git -r -q ./

  echo "File change detected, reloading."

  if [ -n "$PID" ]; then
    # Kill process and restart
    kill $PID > /dev/null
  fi
  sleep 0.2
done
