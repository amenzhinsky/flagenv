workflow "Test" {
  on = "push"
  resolves = ["go test"]
}

action "go test" {
  uses = "docker://golang:1.12-alpine"
  args = ["/bin/sh", "-c", "go test -cover"]
}
