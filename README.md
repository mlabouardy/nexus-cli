[![CircleCI](https://circleci.com/gh/mlabouardy/nexus-cli.svg?style=svg)](https://circleci.com/gh/mlabouardy/nexus-cli) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

<div align="center">
<img src="logo.png" width="60%"/>
</div>

Nexus CLI for Docker Registry

## Usage

<div align="center">
<img src="example.png"/>
</div>

## Download

Below are the available downloads for the latest version of Nexus CLI (1.0.0-beta). Please download the proper package for your operating system and architecture.

### Linux:

```
wget https://s3.eu-west-2.amazonaws.com/nexus-cli/1.0.0-beta/linux/nexus-cli
```

### Windows:

```
wget https://s3.eu-west-2.amazonaws.com/nexus-cli/1.0.0-beta/windows/nexus-cli
```

### Mac OS X:

```
wget https://s3.eu-west-2.amazonaws.com/nexus-cli/1.0.0-beta/osx/nexus-cli
```

### OpenBSD:

```
wget https://s3.eu-west-2.amazonaws.com/nexus-cli/1.0.0-beta/openbsd/nexus-cli
```

### FreeBSD:

```
wget https://s3.eu-west-2.amazonaws.com/nexus-cli/1.0.0-beta/freebsd/nexus-cli
```

## Installation

To install the library and command line program, use the following:

```
go get -u github.com/mlabouardy/nexus-cli
```

## Available Commands

```
$ nexus-cli configure
```

```
$ nexus-cli image ls
```

```
$ nexus-cli image tags -name mlabouardy/nginx
```

```
$ nexus-cli image info -name mlabouardy/nginx -tag 1.2.0
```

```
$ nexus-cli image delete -name mlabouardy/nginx -tag 1.2.0
```

```
$ nexus-cli image tags -name mlabouardy/nginx -v -e v1 -e latest
```

```
$ nexus-cli image tags -name mlabouardy/nginx -e 1.2
```

```
$ nexus-cli image delete -name mlabouardy/nginx -e '!feature'
```

```
$ nexus-cli image delete -name mlabouardy/nginx -keep 4
```

## Caveats

Deletion of image tags is done using a tag name, but rather using a checksum of the image. If you push an image to a registry more than once (e.g. as `1.0.0` and also as `latest`). The deletion will still use the image checksum. Thus the deletion of a single tag is no problem. If the tag is not unique, the deletion will **delete a random tag** matching the checksum.

## Tutorials

* [Cleanup old Docker images from Nexus Repository](http://www.blog.labouardy.com/cleanup-old-docker-images-from-nexus-repository/)
