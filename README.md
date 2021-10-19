<h1 align="center">kubectl-images</h1>
<p align="center">
  <em>🕸 Show container images used in the cluster.</em>
</p>

kubectl-images makes use of the `kubectl` command. It first calls `kubectl get pods` to retrieve pods details and filters out the container image information of each pod, then prints out the final result in a table view.

### 🔰 Installation

Krew
```shell
$ kubectl krew install images
Updated the local copy of plugin index.
Installing plugin: images
Installed plugin: images
\
 | Use this plugin:
 | 	kubectl images
 | Documentation:
 | 	https://github.com/chenjiandongx/kubectl-images
/
```

Build from source code
```shell
$ git clone https://github.com/chenjiandongx/kubectl-images.git
$ cd kubectl-images && go build -ldflags="-s -w" -o kubectl-images . && mv ./kubectl-images /usr/local/bin
$ kubectl images --help
```

Download the binary
```shell
# Refer to the link: https://github.com/chenjiandongx/kubectl-images/releases
# Download the binary and then...
$ chmod +x kubectl-images && mv kubectl-images /usr/local/bin/
$ kubectl images --help
```

### 📝 Usage

```shell
~ 🐶 kubectl images --help
Show container images used in the cluster.

Usage:
  kubectl-images [podname-regex] [flags]

Examples:
  # display a table of all images in current namespace using podName/containerName/containerImage as columns.
  kubectl images

  # display a table of images that match 'nginx' podname regex in 'dev' namespace using podName/containerImage as columns.
  kubectl images -n dev nginx -c 1,2

Flags:
  -A, --all-namespaces         if present, list images in all namespaces.
  -c, --columns string         specify the columns to display, separated by comma. [0:Namespace, 1:PodName, 2:ContainerName, 3:ContainerImage] (default "1,2,3")
      --context string         The name of the kubeconfig context to use.
  -h, --help                   help for kubectl-images
  -k, --kubeconfig string      path to the kubeconfig file to use for CLI requests.
  -n, --namespace string       if present, list images in the specified namespace only. Use current namespace as fallback.
  -o, --output-format string   output format. [json|table] (default "table")
      --version                version for kubectl-images
```

### 🔖 Glances

![image](https://user-images.githubusercontent.com/19553554/74729593-a9201e00-527f-11ea-8325-a4c332dde783.png)
![image](https://user-images.githubusercontent.com/19553554/74729607-ade4d200-527f-11ea-938d-892158d7560f.png)

### 📃 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
