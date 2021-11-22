#!/usr/bin/env sh

readonly session="nomad-local"

tmux switch-client -t "$session"
tmux select-pane -t 0
tmux send-keys C-c

tmux kill-session -t "$session" 2> /dev/null || true
