# sandbox

## Install k3d

```sh
wget -q -O - https://raw.githubusercontent.com/rancher/k3d/master/install.sh | bash
```

### Using a local registry

[docs](https://github.com/rancher/k3d/blob/master/docs/registries.md#using-a-local-registry)
```sh
k3d create --enable-registry -w 5
```

