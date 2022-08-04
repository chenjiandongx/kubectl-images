<h1 align="center">kubectl-images</h1>
<p align="center">
  <em>🕸 Show container images used in the cluster.</em>
</p>

kubectl-images makes use of the `kubectl` command. It first calls `kubectl get pods` to retrieve pods details and
filters out the container image information of each pod, then prints out the final result in a table/json/yaml view.

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

  # display images info in yaml format
  kubectl images -oy

  # display a table of images that match 'nginx' podname regex in 'dev' namespace using podName/containerImage as columns.
  kubectl images -n dev nginx -c 1,2

Flags:
  -A, --all-namespaces         if present, list images in all namespaces.
  -c, --columns string         specify the columns to display, separated by comma. [0:Namespace, 1:PodName, 2:ContainerName, 3:ContainerImage, 4:ImagePullPolicy] (default "1,2,3")
  -C, --context string         The name of the kubeconfig context to use.
  -h, --help                   help for kubectl-images
  -k, --kubeconfig string      path to the kubeconfig file to use for CLI requests.
  -n, --namespace string       if present, list images in the specified namespace only. Use current namespace as fallback.
  -o, --output-format string   output format. [json(j)|table(t)|yaml(y)] (default "table")
  -u, --unique                 Unique images group by namespace/container/images/pullPolicy.
      --version                version for kubectl-images
```

### 🔖 Glances

```shell
~ 🐶 kubectl images -n kube-system -oy dns
- pod: coredns-78fcd69978-9pbjh
  container: coredns
  image: k8s.gcr.io/coredns/coredns:v1.8.4
- pod: coredns-78fcd69978-jh7m2
  container: coredns
  image: k8s.gcr.io/coredns/coredns:v1.8.4

~ 🐶 kubectl images -A -c 0,1,3
[Summary]: 2 namespaces, 11 pods, 11 containers and 9 different images
+-------------+----------------------------------------+--------------------------------------------+
|  Namespace  |                  Pod                   |                   Image                    |
+-------------+----------------------------------------+--------------------------------------------+
| kube-system | coredns-78fcd69978-9pbjh               | k8s.gcr.io/coredns/coredns:v1.8.4          |
+             +----------------------------------------+                                            +
|             | coredns-78fcd69978-jh7m2               |                                            |
+             +----------------------------------------+--------------------------------------------+
|             | etcd-docker-desktop                    | k8s.gcr.io/etcd:3.5.0-0                    |
+             +----------------------------------------+--------------------------------------------+
|             | kube-apiserver-docker-desktop          | k8s.gcr.io/kube-apiserver:v1.22.5          |
+             +----------------------------------------+--------------------------------------------+
|             | kube-controller-manager-docker-desktop | k8s.gcr.io/kube-controller-manager:v1.22.5 |
+             +----------------------------------------+--------------------------------------------+
|             | kube-proxy-vc7fv                       | k8s.gcr.io/kube-proxy:v1.22.5              |
+             +----------------------------------------+--------------------------------------------+
|             | kube-scheduler-docker-desktop          | k8s.gcr.io/kube-scheduler:v1.22.5          |
+             +----------------------------------------+--------------------------------------------+
|             | storage-provisioner                    | docker/desktop-storage-provisioner:v2.0    |
+             +----------------------------------------+--------------------------------------------+
|             | vpnkit-controller                      | docker/desktop-vpnkit-controller:v2.0      |
+-------------+----------------------------------------+--------------------------------------------+
| nginx       | nginx-deployment-66b6c48dd5-s9wv5      | nginx:1.14.2                               |
+             +----------------------------------------+                                            +
|             | nginx-deployment-66b6c48dd5-wmn9x      |                                            |
+-------------+----------------------------------------+--------------------------------------------+

~ 🐶 kubectl images -A -c 0,1,3 -u
[Summary]: 2 namespaces, 11 pods, 11 containers and 9 different images
+-------------+----------------------------------------+--------------------------------------------+
|  Namespace  |                  Pod                   |                   Image                    |
+-------------+----------------------------------------+--------------------------------------------+
| kube-system | coredns-78fcd69978-9pbjh               | k8s.gcr.io/coredns/coredns:v1.8.4          |                                      +
+             +----------------------------------------+--------------------------------------------+
|             | etcd-docker-desktop                    | k8s.gcr.io/etcd:3.5.0-0                    |
+             +----------------------------------------+--------------------------------------------+
|             | kube-apiserver-docker-desktop          | k8s.gcr.io/kube-apiserver:v1.22.5          |
+             +----------------------------------------+--------------------------------------------+
|             | kube-controller-manager-docker-desktop | k8s.gcr.io/kube-controller-manager:v1.22.5 |
+             +----------------------------------------+--------------------------------------------+
|             | kube-proxy-vc7fv                       | k8s.gcr.io/kube-proxy:v1.22.5              |
+             +----------------------------------------+--------------------------------------------+
|             | kube-scheduler-docker-desktop          | k8s.gcr.io/kube-scheduler:v1.22.5          |
+             +----------------------------------------+--------------------------------------------+
|             | storage-provisioner                    | docker/desktop-storage-provisioner:v2.0    |
+             +----------------------------------------+--------------------------------------------+
|             | vpnkit-controller                      | docker/desktop-vpnkit-controller:v2.0      |
+-------------+----------------------------------------+--------------------------------------------+
| nginx       | nginx-deployment-66b6c48dd5-s9wv5      | nginx:1.14.2                               |
+-------------+----------------------------------------+--------------------------------------------+

~ 🐶 kubectl images -c 0,1,2,3,4 -n nginx -oj
[
 {
  "namespace": "nginx",
  "pod": "nginx-deployment-66b6c48dd5-s9wv5",
  "container": "nginx",
  "image": "nginx:latest",
  "imagePullPolicy": "IfNotPresent"
 },
 {
  "namespace": "nginx",
  "pod": "nginx-deployment-66b6c48dd5-wmn9x",
  "container": "nginx",
  "image": "nginx:latest",
  "imagePullPolicy": "IfNotPresent"
 }
]
```

### 📃 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
