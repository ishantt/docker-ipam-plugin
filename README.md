# docker-ipam-plugin
A sample docker IPAM plugin

This plugin is created just to demonstrate a example of Docker IPAM Plugin.


## Installation and Usage
### 1) Run the plugin as Docker Managed Plugin (V2 Plugin)
```
git clone http://github.com/ishantt/docker-ipam-plugin
cd docker-ipam-plugin
make all
```

This will install and enable the docker plugin.

### 2) Check the plugin status 
```
$ docker plugin ls
ID                  NAME                  DESCRIPTION                     ENABLED
b42c13fbaea9        ishant8/sdip:latest   Sample IPAM plugin for Docker   true
```

### 3) Test the plugin by creating a network

Create a network with the plugin
```
$ docker network create --ipam-driver ishant8/sdip:latest test
```
Check your syslog for results

### 4) Use this network to launch a container
```
$ docker run -it --network net1 alpine sh
/ # ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
69: eth0@if70: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether 02:42:3d:38:b8:9d brd ff:ff:ff:ff:ff:ff
    inet 192.168.10.3/24 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::42:3dff:fe38:b89d/64 scope link
       valid_lft forever preferred_lft forever
/ #
```

## References
1) https://github.com/vieux/docker-volume-sshfs
2) https://docs.docker.com/engine/extend/
3) https://github.com/docker/libnetwork/blob/master/docs/ipam.md
