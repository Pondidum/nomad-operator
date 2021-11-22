#!/usr/bin/env sh

readonly session="nomad-local"

tmux kill-session -t "$session" 2> /dev/null || true
tmux start-server
tmux new-session -d -s "$session"

tmux set-layout even-vertical

tmux select-pane -t 0
tmux send-keys "nomad agent -dev" C-m

tmux splitw -v
tmux select-pane -t 1
tmux send-keys "cd operator; go build; while ! nomad status; do sleep 1; done; ./operator" C-m

tmux splitw -v
tmux select-pane -t 2
tmux send-keys "nomad job run example.nomad"

tmux select-pane -t 2
tmux attach-session -t $session
