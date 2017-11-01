# nexus-cli
Nexus CLI for Docker Registry v2

<img src="example.png"/>

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
$ nexus-cli image delete -name mlabouardy/nginx -keep 4
```

```
$ nexus-cli image delete -keep 4
```
