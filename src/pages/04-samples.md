# Samples
---
## Projects with Docker
```yaml
- group: [docker, dkr]
  name: [start, s]
  cmd: docker-compose -p myproject -f docker-compose.yml up -d
  cwd: .
  desc: Start/restart docker containers

# No need to rewrite the alias
- group: docker
  name: [stop, k]
  cmd: docker-compose -p myproject -f docker-compose.yml down --remove-orphans
  cwd: .
  desc: Stop all docker containers

- group: docker
  name: [build, b]
  cmd: docker-compose -p myproject -f docker-compose.yml build{{if not .cache}} --no-cache{{end}}
  cwd: .
  desc: Rebuild the docker container images
  flags:
    - name: cache
      desc: Enable/disable cache
      default: true
      test: $bool

- group: docker
  name: [logs, l]
  cmd: docker logs --tail {{.lines}} -f {{.container}}
  args:
    - name: container
      desc: Container name or id
      default: myproject_web_1
      required: false
  flags:
    - name: lines
      desc: Number of lines to load
      default: 20
      test: $uint
  desc: Load last lines of server logs
```

Submit your samples with a push request at https://github.com/blunt1337/wtfcmd/blob/gh-pages/src/pages/01-samples.md.