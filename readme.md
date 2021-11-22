# Nomad Operator Example

Repostiory to go along with my [The Operator Pattern in Nomad](https://andydote.co.uk/2021/11/22/nomad-operator-pattern/) blog post.

## Usage

If you have tmux installed, you can run `start.sh` to start the demo; it will start nomad, build the operator, and start it for you, and give you the next command to run the demo app.  Exit and cleanup by running `stop.sh`

Otherwise:

1. start Nomad locally
  ```bash
  nomad agent -dev
  ```
2. build and run the operator
  ```bash
  cd operator
  go build
  ./operator
  ```
3. Register an example application to nomad
  ```bash
  nomad job run example.nomad
  ```
4. Open a browser to http://localhost:4646 to see the jobs in the UI
