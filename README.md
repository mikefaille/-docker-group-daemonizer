# docker-group-daemonizer
Fast and useful script to daemonize docker groups

## Requirements 
- Linux 3.10+
- Upstart or systemd (for now, systemd must be activated from the code)
- docker 1.8+
- To build it: Golang, GCC and linux-headers are needed as the Netlink Linux interface is used


## Build 
  
  ```bash
  go get github.com/mikefaille/docker-group-daemonizer
  cd $GOPATH/src/github.com/mikefaille/docker-group-daemonizer
  go build .
  # docker-group-daemonizer bin is now generated
  ```

## Usage 
  Create a Docker group using the following naming scheme : `docker-eqX`
  where X is a number between 1 and 100. After exec it : `sudo docker-group-daemonizer`
 
## Result 
Each docker-eqx group having users will obtain dedicated/isolated docker daemon.

To use these daemons with the Docker CLI, we need to specify appropriate DOCKER_HOST env. variable for each user.

DOCKER_HOST must contain docker-eqx socket named /var/run/docker-eqx.sock
## Example 

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
## Labels
 [equipe=1]
```

## Troubleshoot

Error no1 : 

```
$ docker pull 
Using default tag: latest
Cannot connect to the Docker daemon. Is 'docker -d' running on this host?
```

Explanation no1 :

You can connect to the socket `unix:///var/run/docker-eq1.sock` but the associated daemon is down. Why ? The daemon's bridge could be down or the bridge didn't receive is IP address. You can check if bridge name exist `docker-eq1` with is ip using following command to list network interfaces : `ip adddr`

Fix no1. :

Just recreate bridges named docker-eqx using `docker-group-daemonizer`.


## TODO

(easy to implement) Since docker-group-daemonizer regenerate all bridge named docker-eqx, each execution stop all associated daemon. We must having a way to add only missing bridge. And, --> optionnaly, reset bridge. 

(easy to implement) Manually specify only one docker's bridge interface at time corresponding docker-eqx to be added.
