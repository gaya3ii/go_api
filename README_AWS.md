# AWS Deployment — ECR + EC2

## Overview

The app is containerized, pushed to Amazon ECR, and run on an EC2 instance.

---

## Prerequisites

- AWS CLI installed and configured (`aws configure`)
- Docker installed locally and on the EC2 instance
- An ECR repository created
- An EC2 instance running (Amazon Linux 2 or Ubuntu)

---

## 1. Authenticate Docker to ECR

```bash
aws ecr get-login-password --region <region> | \
  docker login --username AWS --password-stdin \
  <account_id>.dkr.ecr.<region>.amazonaws.com
```

---

## 2. Build the Docker image

```bash
docker build -t go-api .
```

---

## 3. Tag the image for ECR

```bash
docker tag go-api:latest \
  <account_id>.dkr.ecr.<region>.amazonaws.com/go-api:latest
```

---

## 4. Push to ECR

```bash
docker push \
  <account_id>.dkr.ecr.<region>.amazonaws.com/go-api:latest
```

---

## 5. On the EC2 instance — authenticate and pull

SSH into your EC2 instance, then:

```bash
# Authenticate
aws ecr get-login-password --region <region> | \
  docker login --username AWS --password-stdin \
  <account_id>.dkr.ecr.<region>.amazonaws.com

# Pull the image
docker pull <account_id>.dkr.ecr.<region>.amazonaws.com/go-api:latest
```

---

## 6. Run the container on EC2

```bash
docker run -d -p 8080:8080 \
  -e DB_HOST=<rds_or_postgres_host> \
  -e DB_USER=postgres \
  -e DB_PASSWORD=<password> \
  -e DB_NAME=goapi \
  -e DB_PORT=5432 \
  <account_id>.dkr.ecr.<region>.amazonaws.com/go-api:latest
```

---

## 7. EC2 Security Group

Port 8080 must be open for inbound traffic:

| Type       | Protocol | Port | Source    |
|------------|----------|------|-----------|
| Custom TCP | TCP      | 8080 | 0.0.0.0/0 |

To update: **EC2 Console → Security Groups → Inbound rules → Add rule**

---

## Verify

```bash
curl http://<ec2_public_ip>:8080/users
```

---

## Useful commands

```bash
# Check running containers on EC2
docker ps

# View container logs
docker logs <container_id>

# Stop the container
docker stop <container_id>

# Pull and restart with latest image
docker pull <account_id>.dkr.ecr.<region>.amazonaws.com/go-api:latest
docker stop <container_id>
docker run -d -p 8080:8080 ... <image>
```
