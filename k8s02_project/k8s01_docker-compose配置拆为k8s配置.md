# docker-compose配置拆为k8s配置



你这个 `docker-compose` 本质上是 **两个服务（HTTP + cron）**，迁移到 Kubernetes（你用 Minikube + kubectl）时，需要拆成几类资源：

👉 对应关系先说清楚（很关键）：

* service（docker-compose）→ Deployment（K8s）
* ports → Service（K8s暴露）
* volumes（配置）→ ConfigMap
* healthcheck → livenessProbe / readinessProbe
* restart → Deployment 自带
* container_name → 不需要（K8s自动管理）

---

# 一、ConfigMap（配置文件）

把 `config.yaml` 挂进去：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gin-kubelet-config
data:
  config.yaml: |
    # 把你的 config.yaml 内容粘到这里
```

---

# 二、HTTP 服务（Deployment）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gin-kubelet-http
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gin-kubelet-http
  template:
    metadata:
      labels:
        app: gin-kubelet-http
    spec:
      containers:
        - name: gin-kubelet-http
          image: gin_kubelet:latest   # 你本地 build 的镜像
          imagePullPolicy: IfNotPresent
          command: ["/app/gin_kubelet", "http", "-c", "/app/config/config.yaml"]
          ports:
            - containerPort: 8080
          volumeMounts:
            - name: config-volume
              mountPath: /app/config/config.yaml
              subPath: config.yaml
            - name: logs
              mountPath: /app/logs

          # 对应 docker-compose healthcheck
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 30
            timeoutSeconds: 10
            failureThreshold: 3

          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10

      volumes:
        - name: config-volume
          configMap:
            name: gin-kubelet-config
        - name: logs
          emptyDir: {}
```

---

# 三、HTTP 对外暴露（Service）

```yaml
apiVersion: v1
kind: Service
metadata:
  name: gin-kubelet-http
spec:
  type: NodePort   # minikube 推荐
  selector:
    app: gin-kubelet-http
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 30080
```

👉 访问方式：

```bash
minikube service gin-kubelet-http
```

---

# 四、Cron 服务（Deployment）

⚠️ 重点：你现在这个 cron 服务是**常驻进程**
所以用 Deployment（不是 K8s CronJob）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gin-kubelet-cron
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gin-kubelet-cron
  template:
    metadata:
      labels:
        app: gin-kubelet-cron
    spec:
      containers:
        - name: gin-kubelet-cron
          image: gin_kubelet:latest
          imagePullPolicy: IfNotPresent
          command: ["/app/gin_kubelet", "cron", "-c", "/app/config/config.yaml"]
          volumeMounts:
            - name: config-volume
              mountPath: /app/config/config.yaml
              subPath: config.yaml
            - name: logs
              mountPath: /app/logs

      volumes:
        - name: config-volume
          configMap:
            name: gin-kubelet-config
        - name: logs
          emptyDir: {}
```

---

# 五、在 Minikube 上部署步骤

### 1️⃣ 启动

```bash
minikube start
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

### 3️⃣ 应用 YAML

```bash
kubectl apply -f configmap.yaml
kubectl apply -f http-deployment.yaml
kubectl apply -f http-service.yaml
kubectl apply -f cron-deployment.yaml
```

---

### 4️⃣ 查看状态

```bash
kubectl get pods
kubectl get svc
kubectl logs -f gin-kubelet-http-xxx
```

---

# 六、几个关键差异总结（帮你避免踩坑）

### 1️⃣ healthcheck → probe

docker:

```yaml
wget http://localhost:8080/health
```

k8s：

```yaml
httpGet:
  path: /health
  port: 8080
```

👉 不需要 wget

---

### 2️⃣ volumes 映射

docker：

```yaml
./config.yaml:/app/config.yaml
```

k8s：
👉 用 ConfigMap + subPath

---

### 3️⃣ 日志目录

你现在：

```yaml
./logs
```

k8s：

```yaml
emptyDir
```

👉 如果你要持久化 → 用 PVC（可以再帮你补）

---

### 4️⃣ container_name

👉 K8s **不允许固定容器名**（Pod 会动态生成）

---

# 如果你想再进阶

可以继续优化：

* ✅ 使用 Ingress（代替 NodePort）
* ✅ 使用 PVC 持久化 logs
* ✅ 把 cron 改成 K8s CronJob（更正宗）
* ✅ 加资源限制（CPU / 内存）

---

如果你愿意，我可以帮你再升级一版：
👉 **生产级 K8s YAML（带 Ingress + HPA + PVC）**
