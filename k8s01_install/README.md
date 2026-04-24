# k8s01_install


1. [Install Docker Desktop on Mac](https://docs.docker.com/desktop/setup/install/mac-install/)
2. [在 macOS 系统上安装和设置 kubectl](https://kubernetes.io/zh-cn/docs/tasks/tools/install-kubectl-macos/)
3. [在 macOS 上启动 minikube ](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download)
4. 


```markdown

kubectl 是“遥控器”，minikube 是“本地小型 Kubernetes 集群”。
🟢 kubectl = “控制工具”，用来操作集群，但它不负责提供集群服务
    kubectl get pods
    kubectl apply -f xxx.yaml
    kubectl logs pod
    kubectl delete deployment

🟡 minikube = “本地 Kubernetes 环境”，用来启动集群，在你电脑上跑一个“单机版 Kubernetes”
    minikube start

我们不能直接操作 minikube 集群，因此要使用一个操作 minikube 集群的工具 kubectl，kubectl 来操作 minikube 启动的集群

| 对比项      | kubectl     | minikube         |
| -------- | ----------- | ---------------- |
| 本质       | 命令行工具       | 本地 Kubernetes 集群 |
| 作用       | 操作 K8s 集群   | 创建/运行 K8s 集群     |
| 是否包含 K8s | ❌ 不包含       | ✅ 包含             |
| 是否独立运行   | ❌ 不行（必须连集群） | ✅ 可以独立运行         |
| 类比       | 遥控器 / CLI   | 虚拟机 / 小型服务器      |
| 使用场景     | 管理任何 K8s 集群 | 本地开发/学习          |

```