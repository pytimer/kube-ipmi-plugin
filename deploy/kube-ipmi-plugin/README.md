# kube-ipmi-plugin

## Installing the Chart

To install the chart with the release name `my-release`:

`$ helm install --name my-release ./kube-ipmi-plugin`

## Uninstalling the Chart

To uninstall/delete the `my-release`:

`$ helm delete my-release --purge`

## Configuration

## Custom image and namespace `my-namespace`:

`$ helm install --name my-release --namespace=my-namespace --set image.repository=<custom-image-name> ./kube-ipmi-plugin`
