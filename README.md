
# Docker Volumes Provisioner
hub.docker.com [viniokil/volumes-provisioner](https://hub.docker.com/r/viniokil/volumes-provisioner)


An golang based simple helper image to essentially make directory
with appropriate owner and permissions to be used by other containers.

Most images require data folders to be under certain permissions,
for mounted volumes directories have to be created and modified
before container start.
e.g.

- [mysql](https://hub.docker.com/_/mysql) data folder needs to be 1000:1000
- [prometheus](https://hub.docker.com/r/prom/prometheus) needs 65534:65534
- [grafana](https://hub.docker.com/r/grafana/grafana) needs 472:472


### Config
Pass config as env var `PROVISION_DIRECTORIES`

uid:gid:mode:dir(s) e.g. 1000:1000:0755:/var/www

Can contain multiple directories separated by space and/or multiple configs separated by ;

e.g.

- `1000:1000:0755:/var/www/html`
- `1000:1000:0755:/var/www /var/www/html`
- `1000:1000:0755:/var/www /var/www/html;1001:1001:0755:/var/app /var/app/cache`

### Local Run example
```
docker run --rm \
-e PROVISION_DIRECTORIES=1000:1000:0755:/var/www \
-v `pwd`/data:/var/www \
viniokil/volumes-provisioner
ls -dhn data
```
### Loging
#### Log level
      # # (Default: ERROR) DEBUG, INFO, WARN, ERROR, FATAL, PANIC

```yaml
version: '3'
services:
  volumes-provisioner:
    image: viniokil/volumes-provisioner
    environment:
      LOG_LEVEL: "DEBUG"
    volumes:
      - "/var/data:/var/data"
    network_mode: none
```

#### Log format
Available log formats:
- text
- logfmt
- json

```yaml
version: '3'
services:
  volumes-provisioner:
    image: viniokil/volumes-provisioner
    environment:
      LOG_FORMAT: "json"
    volumes:
      - "/var/data:/var/data"
    network_mode: none
```
### Docker Compose prometheus example
Here using depends_on we make sure directory with correct
permissions is created before prometheus starts
```yaml
version: '2'
services:
  volumes-provisioner:
    image: viniokil/volumes-provisioner
    environment:
      PROVISION_DIRECTORIES: "65534:65534:0755:/var/data/prometheus/data"
    volumes:
      - "/var/data:/var/data"
    network_mode: none

  prometheus:
    image: prom/prometheus:v2.3.2
    ports:
      - "9090:9090"
    depends_on:
      - volumes-provisioner
    volumes:
      - "/var/data/prometheus/data:/prometheus/data"

```

### Build locally
```sh
docker build -t volumes-provisioner .
docker run --rm \
-e PROVISION_DIRECTORIES=1000:1000:0755:/var/www \
-v `pwd`/data:/var/www \
volumes-provisioner
ls -dhn data
```

### Notes
Won't work for docker for macos volume mounts
