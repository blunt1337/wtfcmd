# Samples
---

## For Laravel projects
```yaml
# Misc
- name: init
  desc: Initialise the project
  cmd: npm install; composer install
  cwd: .

# Frontend
- group: [frontend, front, f]
  name: [watch, w]
  cmd: npm run watch
  desc: Recompile js/css/html files on source change

- group: frontend
  name: [compile, c]
  cmd: npm run production
  desc: Recompile js/css/html files for production

# Database
- group: db
  name: update
  cmd: php /app/artisan migrate --force
  desc: Run updates on the database

- group: db
  name: rollback
  cmd: php /app/artisan migrate:rollback
  desc: Remove last database updates

- group: db
  name: seed
  cmd: php /app/artisan db:seed
  desc: Add test records in the database
```

## For docker-compose projects, named myproject
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
