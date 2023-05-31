# demo
Demonstration of the cluster that have three environments(main, feature-1, feature-2) and that is deployed application with Services A, B, C (Web servers) and Service MySQL.

# Usage
## Create Resources
Please follow steps to create there resources.
- Proxy Controller
- Proxy A/B/C
- Service A[main, feature-1]
- Service B[main, feature-2]
- Service C[main]
- Service MySQL[main, feature-1, feature-2]
- Proxy Controller DB
### Docker Compose with One Host
One host has all resources.
```
cd docker/all-in-one
docker compose up 
```

### Docker Compose with Three Hosts(Separated Web Servers, a MySQL Server, and Proxy Controller DB)
#### Host Web
It has the following resources.
- Proxy Controller
- Proxy A/B/C
- Service A[main, feature-1]
- Service B[main, feature-2]
- Service C[main]
```
cd docker/without-db
vim docker-compose.yaml # Please replace the IP address 192.168.0.3 with that of the Host MySQL and 192.168.0.4 with that of the Host Controller DB.
docker compose up
```

#### Host MySQL
It has Service MySQL[main, feature-1, feature-2].
```
cd docker/only-service-mysql
docker compose up
```
Please confirm that the Host Web can access the port 13306, 13307, and 13308 of the Host MySQL.

#### Host Controller DB
It has Proxy Controller DB.
1. Install MySQL
2. Create MySQL user `picop` with password `picop`
3. Configure it to allow the Host Web to login `picop`
4. Create database `picop`
5. Create tables referring to `proxy-controller/db-schema`

Please confirm that the Host Web can access the port 3306 of the Host Controller DB.

### Kubernetes
Please prepare Kubernetes cluster and the other two hosts.

#### Kubernetes Cluster
It has the following resources.
- Proxy Controller
- Proxy A/B/C
- Service A[main, feature-1]
- Service B[main, feature-2]
- Service C[main]
```
cd kubernetes/manifests
vim proxy-mysql.yaml # Please replace the IP address 192.168.0.3 with that of the Host MySQL.
vim proxy-controller.yaml # Please replace the IP address 192.168.0.4 with that of the Host Controller DB.
kubectl apply namespace.yaml
kubectl apply proxy-controller.yaml
kubectl apply proxy-a.yaml
kubectl apply proxy-b.yaml
kubectl apply proxy-c.yaml
kubectl apply proxy-mysql.yaml
./script/create-service-a.sh main | kubectl apply -f -
./script/create-service-b.sh main | kubectl apply -f -
./script/create-service-c.sh main | kubectl apply -f -
./script/create-service-a.sh feature-1 | kubectl apply -f -
./script/create-service-b.sh feature-2 | kubectl apply -f -
```

#### Host MySQL and Host Controller DB
It is the same as the case of Docker Compose with Three Hosts.

## Configure Proxy Controller
Please get Proxy Controller URL. If Docker Compose, it is `http://<the Host IP address with Proxy Controller>:8080`. If Kubernetes, it is `http://<k8s cluster IP address>:31001`.

```
cs script
./register-proxies.sh <Proxy Controller URL>
vim ./register-routes.yaml # Please replace the IP address 192.168.0.3 with that of the Host MySQL.
./register-routes.sh <Proxy Controller URL>
```

## Send Requests with PiCoP
You can use [picop-curl](https://github.com/picop-rd/picop-curl).
Please get Proxy A URL. If Docker Compose, it is `http://<the Host IP address with Proxy A>:9001`. If Kubernetes, it is `http://<k8s cluster IP address>:31002`.

Environment `main`(or other environments such as `feature-100`)
```
$ picop-curl --env-id main --url <Proxy A URL> --method POST --data test-main
This is service-a-main
This is service-b-main
This is service-c-main

$ picop-curl --env-id main --url <Proxy A URL> --method GET
This is service-a-main
This is service-b-main
This is service-c-main
data{ id: 0, content: test-main }

$ picop-curl --env-id feature-100 --url <Proxy A URL> --method POST --data test-feature-100
This is service-a-main
This is service-b-main
This is service-c-main

$ picop-curl --env-id feature-100 --url <Proxy A URL> --method GET
This is service-a-main
This is service-b-main
This is service-c-main
data{ id: 0, content: test-main }
data{ id: 1, content: test-feature-100 }

```

Environment `feature-1`
```
$ picop-curl --env-id feature-1 --url <Proxy A URL> --method POST --data test-feature-1
This is service-a-feature-1
This is service-b-main
This is service-c-main

$ picop-curl --env-id feature-1 --url <Proxy A URL> --method GET
This is service-a-feature-1
This is service-b-main
This is service-c-main
data{ id: 0, content: test-feature-1 }

```

Environment `feature-2`
```
$ picop-curl --env-id feature-2 --url <Proxy A URL> --method POST --data test-feature-2
This is service-a-main
This is service-b-feature-2
This is service-c-main

$ picop-curl --env-id feature-2 --url <Proxy A URL> --method GET
This is service-a-main
This is service-b-feature-2
This is service-c-main
data{ id: 0, content: test-feature-2 }

```


