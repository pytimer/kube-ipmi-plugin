# kube-ipmi-plugin
The plugin collect ipmi info and report the Kubernetes nodes

## Development

Go 1.12+

### Build binary

`make build`

### Build Docker image

`make docker-build`

If you want to push docker image, you can use `make docker-push` command.


## Annotations

The `kube-ipmi-plugin` collect the ipmi info and save the information to the Kubernetes nodes annotation.

The annotation key is: `ipmi.alpha.kubernetes.io/net`.

#### Example:

```json
ipmi.alpha.kubernetes.io/net:{"default_gateway_ip":"192.168.254.254","default_gateway_mac":"12:34:56:78:ab:cd","ip_address":"192.168.254.10","ip_address_source":"Static Address","mac_address":"ab:cd:ef:gh:12:34","subnet_mask":"255.255.255.0"}
```

Use `kubectl describe node <nodeName>`, you can lookup this node ipmi network ip address, gateway, etc.

```sh
Annotations:        ipmi.alpha.kubernetes.io/net={"default_gateway_ip":"192.168.254.254","default_gateway_mac":"12:34:56:78:ab:cd","ip_address":"192.168.254.10","ip_address_source":"Static Address","mac_address":"ab:cd:...
```

## Run

### Running locally using Binary

You can get the source code, build the binary. Or use `go install`

```sh
$ go get github.com/pytimer/kube-ipmi-plugin
$ cd $GOPATH/src/github.com/pytimer/kube-ipmi-plugin
$ make build
```

```sh
$ go install github.com/pytimer/kube-ipmi-plugin
```

Make sure the plugin install in `$GOPATH/bin`

### Running locally using Docker

Pull the `kube-ipmi-plugin` Docker image and run it. The latest image is `pytimer/kube-ipmi-plugin`, if you want to specified version, you can search this plugin images on dockerhub.

```sh
$ docker run -it --rm --network=host --device=/dev/ipmi0 pytimer/kube-ipmi-plugin
$ docker run -it --rm --network=host --privileged -v /dev/ipmi0:/dev/ipmi0 pytimer/kube-ipmi-plugin
```

### Running the Kubernetes cluster

You can use `helm` to install this plugin if your Kubernetes cluster support helm. If not, you can use Kubernetes manifests to install.

- [charts](deploy/kube-ipmi-plugin)

- [manifests](deploy/manifests)



