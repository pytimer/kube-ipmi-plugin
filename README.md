# kube-ipmi-plugin
The plugin collect ipmi info and report the Kubernetes nodes

## Development

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
