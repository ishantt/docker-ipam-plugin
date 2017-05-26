# docker-ipam-plugin
A sample docker IPAM plugin

This plugin is created just to demonstate a example of Docker IPAM Plugin.


## Installation and Usage
### 1) To Run the plugin as a external container (legacy plugin)

```
git clone http://github.com/ishantt/docker-ipam-plugin
cd docker-ipam-plugin
make build-image
```

This will build a image with name ishant8/sdip:rootfs

To run the plugin container execute

```
docker run -v /run/docker:/run/docker ishant8/sdip:rootfs
```

Create a network with the IPAM Driver 

```
docker network create --ipam-driver sdip net1
```

Use this network to launch a container
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


### 2) To Run the plugin as Docker Managed Plugin (V2 Plugin)
```
git clone http://github.com/ishantt/docker-ipam-plugin
cd docker-ipam-plugin
make all
```

This will install and enable the docker plugin.

To see the plugin status execute
```
$ docker plugin ls
ID                  NAME                  DESCRIPTION                     ENABLED
b42c13fbaea9        ishant8/sdip:latest   Sample IPAM plugin for Docker   true
```

Create a network with the plugin
```
$ docker network create --ipam-driver ishant8/sdip:latest test
```

Check your syslog for results

This is still a WIP

## References
1) https://github.com/vieux/docker-volume-sshfs
2) https://docs.docker.com/engine/extend/
3) https://github.com/docker/libnetwork/blob/master/docs/ipam.md
