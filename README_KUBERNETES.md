# Kubernetes Deployment (local — minikube)

## Overview

The app runs as a 3-replica Deployment behind a LoadBalancer Service, backed by a
single-replica Postgres Deployment. DB connection settings come from a ConfigMap
and the DB password from a Secret, both injected as environment variables via
`envFrom` — no image rebuild needed to change them.

Manifests live in [k8s/](k8s/):

| File               | Kind                  | Purpose                                   |
|--------------------|-----------------------|--------------------------------------------|
| `configmap.yaml`   | ConfigMap             | `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER` |
| `secret.yaml`       | Secret                | `DB_PASSWORD`                             |
| `postgres.yaml`     | Deployment + Service  | Postgres pod, reads its own env from the ConfigMap/Secret above |
| `deployment.yaml`   | Deployment + Service  | go-api pods, `envFrom` the ConfigMap + Secret |

`deployment.yaml` sets `imagePullPolicy: Never`, so the image must exist on the
cluster node already — it's never pulled from a registry. That's why the image
load step below is required after every rebuild.

---

## Prerequisites

- [minikube](https://minikube.sigs.k8s.io/docs/start/) running (`minikube start`)
- `kubectl` pointed at the minikube context
- Docker (to build the image)

---

## 1. Start the cluster

```bash
minikube start
```

---

## 2. Build the image and load it into minikube

minikube runs its own Docker daemon inside the cluster VM/container, separate
from your host Docker. Building normally with `docker build` only puts the
image on your host, so it has to be loaded in explicitly:

```bash
docker build -t go-api:latest .
minikube image load go-api:latest
```

---

## 3. Apply the manifests

```bash
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/deployment.yaml
```

Or all at once:

```bash
kubectl apply -f k8s/
```

---

## 4. Verify

```bash
kubectl get pods
kubectl get svc
```

Wait until `postgres` and all `go-api` pods show `Running` — `go-api` will
crash-loop if it starts before Postgres is ready to accept connections.

---

## 5. Access the service

minikube doesn't provision a real cloud LoadBalancer, so use a tunnel or the
built-in service URL instead of expecting an external IP:

```bash
minikube service go-api-service --url
```

or, in a separate terminal:

```bash
minikube tunnel
```

Then hit the API the same way as [README.md](README.md)'s API Endpoints section, e.g.:

```bash
curl http://<url-from-above>/users
```

---

## ConfigMap & Secret vs. local `.env`

Same idea as [docker-compose.yml](docker-compose.yml) and the local `.env` file — just a
different delivery mechanism. `envFrom.configMapRef` / `secretRef` in
[k8s/deployment.yaml](k8s/deployment.yaml) inject `DB_HOST`, `DB_USER`, `DB_PASSWORD`,
`DB_NAME`, `DB_PORT` as real process environment variables before the app
starts, so `os.Getenv` in [db/db.go](db/db.go) picks them up exactly as it would
locally or under docker-compose — the Go code doesn't know or care which of
the three set them.

**Note:** [k8s/secret.yaml](k8s/secret.yaml) currently has the DB password as
plaintext `stringData` committed to the repo. That's fine for local
minikube practice, but treat it as a placeholder — don't reuse that pattern
if this ever points at a real database.

---

## Useful commands

```bash
# Logs for a specific go-api pod
kubectl logs -f deployment/go-api

# Logs for postgres
kubectl logs -f deployment/postgres

# Describe a pod (debug crash loops, image pull issues, etc.)
kubectl describe pod <pod-name>

# Restart the deployment after loading a new image
minikube image load go-api:latest
kubectl rollout restart deployment/go-api

# Tear everything down
kubectl delete -f k8s/
```
