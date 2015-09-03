# docker-group-daemonizer
Fast and usefull script to daemonize docker groups

Requirement : 
- Linux 3.10+
- Upstart or systemd (for now systemd must activated from code)
- docker 1.8+
- to build it : Golang, GCC and linux-headers is needed since Netlink Linux interface is used


Build : 
  
  ```bash
  git clone https://github.com/mikefaille/docker-group-daemonizer.git
  cd docker-group-daemonizer
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
