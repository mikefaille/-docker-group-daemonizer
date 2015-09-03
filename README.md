# docker-group-daemonizer
Fast and usefull script to daemonize docker groups

Requirement : 
- Linux 3.10+
- Upstart or systemd (for now systemd must activated from code)
- docker 1.8+
- to build it : Golang, GCC and linux-headers is needed since Netlink Linux interface is used


Build : 
  
  ```bash
  go get github.com/mikefaille/docker-group-daemonizer
  cd $GOPATH/src/github.com/mikefaille/docker-group-daemonizer
  go build .
  # docker-group-daemonizer bin is now generated
  ```

Usage : 
  Create docker group using this next naming schemes : `docker-eqX`
  where X is a number between 1 and 100. After exec it : `sudo docker-group-daemonizer`
 
Result : 
Each docker-eqx group having users will obtain dedicated/isolated docker deamon.

To use these deamon using docker cli, we need to specify appropritate DOCKER_HOST env. variable for each user.

DOCKER_HOST must contain docker-eqx socket named /var/run/docker-eqx.sock
Example : 

`$ groups`

michael docker-eq1

`ls -la /var/run/docker-eq1.sock`

srw-rw---- 1 root docker-eq1 0 Sep  3 15:57 /var/run/docker-eq1.sock

```bash
export DOCKER_HOST=unix:///var/run/docker-eq1.sock
docker pull busybox
docker images
```

And, `docker info` should output something like this : 
```
(...)
Labels:
 [equipe=1]
```
