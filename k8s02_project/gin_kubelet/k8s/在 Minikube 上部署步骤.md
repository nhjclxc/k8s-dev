# 在 Minikube 上部署步骤

启动需要的一些镜像
```shell
# 一、最核心镜像（必须有）1️⃣ kicbase（最关键）
gcr.io/k8s-minikube/kicbase:v0.0.48 【新地址：docker.io/kicbase/stable:v0.0.48】
# 二、Kubernetes 核心组件镜像2️⃣ API Server👉 集群入口（kubectl 就是访问它）
k8s.gcr.io/kube-apiserver
# 二、Kubernetes 核心组件镜像3️⃣ Controller Manager👉 控制器（副本数、节点状态等）
k8s.gcr.io/kube-controller-manager
# 二、Kubernetes 核心组件镜像4️⃣ Scheduler👉 负责 Pod 调度
k8s.gcr.io/kube-scheduler
# 二、Kubernetes 核心组件镜像5️⃣ etcd👉 存储所有集群数据
k8s.gcr.io/etcd
# 二、Kubernetes 核心组件镜像6️⃣ CoreDNS👉 集群内 DNS
k8s.gcr.io/coredns
# 二、Kubernetes 核心组件镜像7️⃣ Pause（非常重要）👉 Pod 的“根容器”
k8s.gcr.io/pause
# 三、minikube 自带组件镜像8️⃣ storage-provisioner👉 默认动态存储（PVC 用）
gcr.io/k8s-minikube/storage-provisioner:v5
# 三、minikube 自带组件镜像9️⃣ ingress（可选）【如果使用了命令 minikube addons enable ingress 会拉取下面这个镜像】
k8s.gcr.io/ingress-nginx/controller

```

[国内镜像chinese-opensource-mirror-site ](https://github.com/SUSTech-CRA/chinese-opensource-mirror-site)
[国内镜像public-image-mirror ](https://github.com/DaoCloud/public-image-mirror)



### 1️⃣ 启动 minikube 集群
```
✅ 1. 清理脏环境（必须做）【这个操作会清理docker里面所有镜像和容器】
minikube delete
docker system prune -af

✅ 2. 手动准备 完整控制面镜像（重点）
docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/kicbase:v0.0.48
docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver:v1.34.0
docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager:v1.34.0
docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler:v1.34.0
docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/coredns:v1.11.1
docker pull registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.10
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/registry.k8s.io/etcd:3.5.12-0-linuxarm64

✅ 3. 打 tag（让 minikube 能识别）
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kicbase:v0.0.48  gcr.io/k8s-minikube/kicbase:v0.0.48
docker tag gcr.io/k8s-minikube/kicbase:v0.0.48  docker.io/kicbase/stable:v0.0.48
docker tag gcr.io/k8s-minikube/kicbase:v0.0.48  registry.k8s.io/k8s-minikube/kicbase:v0.0.48
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver:v1.34.0 registry.k8s.io/kube-apiserver:v1.34.0
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager:v1.34.0 registry.k8s.io/kube-controller-manager:v1.34.0
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler:v1.34.0 registry.k8s.io/kube-scheduler:v1.34.0
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/coredns:v1.11.1 registry.k8s.io/coredns:v1.11.1
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.10 registry.k8s.io/pause:3.10
docker tag swr.cn-north-4.myhuaweicloud.com/ddn-k8s/registry.k8s.io/etcd:3.5.12-0-linuxarm64 registry.k8s.io/etcd:3.5.12

✅ 4. 再启动（带资源）
➜  debug git:(main) ✗ minikube start \
  --driver=docker \
  --memory=4096 \
  --image-mirror-country=cn
😄  Darwin 15.1.1 (arm64) 上的 minikube v1.37.0
E0424 11:22:04.909055   47956 start.go:829] api.Load failed for minikube: filestore "minikube": Docker machine "minikube" does not exist. Use "docker-machine ls" to list machines. Use "docker-machine create" to add a new one.
E0424 11:22:04.909166   47956 start.go:829] api.Load failed for minikube: filestore "minikube": Docker machine "minikube" does not exist. Use "docker-machine ls" to list machines. Use "docker-machine create" to add a new one.
✨  根据现有的配置文件使用 docker 驱动程序
👍  在集群中 "minikube" 启动节点 "minikube" primary control-plane
🚜  正在拉取基础镜像 v0.0.48 ...
❗  minikube was unable to download registry.cn-hangzhou.aliyuncs.com/google_containers/kicbase:v0.0.48, but successfully downloaded gcr.io/k8s-minikube/kicbase:v0.0.48 as a fallback image
🔥  创建 docker container（CPU=2，内存=4096MB）...
🐳  正在 Docker 28.4.0 中准备 Kubernetes v1.34.0…
🔗  配置 bridge CNI (Container Networking Interface) ...
🔎  正在验证 Kubernetes 组件...
    ▪ 正在使用镜像 registry.cn-hangzhou.aliyuncs.com/google_containers/storage-provisioner:v5
🌟  启用插件： storage-provisioner, default-storageclass
🏄  完成！kubectl 现在已配置，默认使用"minikube"集群和"default"命名空间
➜  debug git:(main) ✗ docker ps                                               
CONTAINER ID   IMAGE                                 COMMAND                   CREATED              STATUS              PORTS                                                                                                                                  NAMES
7cc692dfb19e   gcr.io/k8s-minikube/kicbase:v0.0.48   "/usr/local/bin/entr…"   About a minute ago   Up About a minute   127.0.0.1:63545->22/tcp, 127.0.0.1:63546->2376/tcp, 127.0.0.1:63548->5000/tcp, 127.0.0.1:63544->8443/tcp, 127.0.0.1:63547->32443/tcp   minikube
➜  debug git:(main) ✗ kubectl get nodes
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   85s   v1.34.0
➜  debug git:(main) ✗ kubectl get pods -A
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE
kube-system   coredns-7ddb67b59b-hgbs2           1/1     Running   0          82s
kube-system   etcd-minikube                      1/1     Running   0          88s
kube-system   kube-apiserver-minikube            1/1     Running   0          88s
kube-system   kube-controller-manager-minikube   1/1     Running   0          88s
kube-system   kube-proxy-whqfj                   1/1     Running   0          82s
kube-system   kube-scheduler-minikube            1/1     Running   0          88s
kube-system   storage-provisioner                1/1     Running   0          87s
➜  debug git:(main) ✗ kubectl cluster-info
Kubernetes control plane is running at https://127.0.0.1:63544
CoreDNS is running at https://127.0.0.1:63544/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.


查询容器是否在运行
➜  ~ docker ps                   
CONTAINER ID   IMAGE                                 COMMAND                   CREATED       STATUS       PORTS                                                                                                                                  NAMES
7cc692dfb19e   gcr.io/k8s-minikube/kicbase:v0.0.48   "/usr/local/bin/entr…"   6 hours ago   Up 6 hours   127.0.0.1:63545->22/tcp, 127.0.0.1:63546->2376/tcp, 127.0.0.1:63548->5000/tcp, 127.0.0.1:63544->8443/tcp, 127.0.0.1:63547->32443/tcp   minikube



如果要打开网页，则执行：minikube dashboard
建议提前准备好以下镜像资源
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kubernetesui/metrics-scraper:v1.0.8
docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kubernetesui/metrics-scraper:v1.0.8  docker.io/kubernetesui/metrics-scraper:v1.0.8
docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kubernetesui/dashboard:v2.7.0
docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/kubernetesui/dashboard:v2.7.0  docker.io/kubernetesui/dashboard:v2.7.0


➜  debug git:(main) ✗ minikube dashboard                                                                                
🤔  正在验证 dashboard 运行情况 ...
🚀  正在启动代理...
🤔  正在验证 proxy 运行状况 ...
🎉  正在使用默认浏览器打开 http://127.0.0.1:64016/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/ ...


```



### 2️⃣ 使用 minikube 的 docker 环境（关键！）

```bash
eval $(minikube docker-env)
```

👉 然后再 build：

```bash
docker build -t gin_kubelet:latest .
```

---

### 3️⃣ 应用 YAML 启动业务服务

如果资源使用了命名空间，那么必须先创建命名空间kubectl apply -f namespace.yaml
```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f http-deployment.yaml
kubectl apply -f http-service.yaml
kubectl apply -f cron-deployment.yaml
```

```
pod重启命令：
kubectl rollout restart deployment gin-kubelet-http -n gin-kubelet-debug

强制重建
kubectl delete pod -l app=gin-kubelet-http -n gin-kubelet-debug

删除pod
kubectl delete pod gin-kubelet-http-6ff7bc4cdc-g9264 -n gin-kubelet-debug
```
---

### 4️⃣ 查看状态

```bash
kubectl get pods
kubectl get svc
kubectl logs -f gin-kubelet-http-xxx
```

---





```shell


➜  debug git:(main) ✗ ls
configmap.yaml       cron-deployment.yaml http-deployment.yaml http-service.yaml    namespace.yaml
➜  debug git:(main) ✗ kubectl get pods         
No resources found in default namespace.
➜  debug git:(main) ✗ kubectl get pods -A
NAMESPACE              NAME                                         READY   STATUS             RESTARTS   AGE
kube-system            coredns-7ddb67b59b-hgbs2                     1/1     Running            0          7m52s
kube-system            etcd-minikube                                1/1     Running            0          7m58s
kube-system            kube-apiserver-minikube                      1/1     Running            0          7m58s
kube-system            kube-controller-manager-minikube             1/1     Running            0          7m58s
kube-system            kube-proxy-whqfj                             1/1     Running            0          7m52s
kube-system            kube-scheduler-minikube                      1/1     Running            0          7m58s
kube-system            storage-provisioner                          1/1     Running            0          7m57s
kubernetes-dashboard   dashboard-metrics-scraper-77bf4d6c4c-m5n2p   1/1     Running            0          3m55s
kubernetes-dashboard   kubernetes-dashboard-855c9754f9-gm8vv        0/1     ImagePullBackOff   0          3m55s
➜  debug git:(main) ✗ kubectl apply -f configmap.yaml
Error from server (NotFound): error when creating "configmap.yaml": namespaces "gin-kubelet-debug" not found
➜  debug git:(main) ✗ 
➜  debug git:(main) ✗ kubectl apply -f namespace.yaml
namespace/gin-kubelet-debug created
➜  debug git:(main) ✗ kubectl apply -f configmap.yaml
configmap/gin-kubelet-config created
➜  debug git:(main) ✗ kubectl apply -f http-deployment.yaml
deployment.apps/gin-kubelet-http created
➜  debug git:(main) ✗ kubectl apply -f http-service.yaml
service/gin-kubelet-http created
➜  debug git:(main) ✗ kubectl apply -f cron-deployment.yaml
deployment.apps/gin-kubelet-cron created



查看pod日志
kubectl logs gin-kubelet-http-6ff7bc4cdc-g9264 -n gin-kubelet-debug
```






