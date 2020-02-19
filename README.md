# sandbox

## Install k3d

```sh
export K3D_INSTALL_DIR="$HOME/bin"
wget -q -O - https://raw.githubusercontent.com/rancher/k3d/master/install.sh | bash
```
or
```sh
brew install k3d
```

### Using a local registry

[docs](https://github.com/rancher/k3d/blob/master/docs/registries.md#using-a-local-registry)
```sh
k3d create --enable-registry
```

