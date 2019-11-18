# HyKuFe
# Hykufe Operator
hykufe.com/hykufe-operator

## Build and run the operator
operator를 동작시키기 전에 CRD가 쿠버네티스 apiserver에 등록되어야 한다.
```
$ kubectl create -f deploy/crds/hykufe-v1alpha1_horovodjob_crd.yaml
```

도커 이미지로 빌드 하고 도커 허브에 저장할 것이다.
### build operator
```shell script
$ operator-sdk generate k8s
$ operator-sdk generate openapi      
$ operator-sdk build
$ go mod vendor

$ opreator-sdk build [도커 허브 ID]/hykuge-operator:v0.0.1
$ sed -i 's|REPLACE_IMAGE|quay.io/example/[도커 허브 ID]/hykuge-operator:v0.0.1|g' deploy/operator.yaml
$ docker push [도커 허브 ID]/hykuge-operator:v0.0.1

```

### Run Operator in kube cluster
```shell script
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/operator.yaml
```

### 확인
```shell script
$ kubectl get deployment
```

## Custrom Resource Definition(CRD)
### HorovodJob


- name
- type
- HorovodJobSpec
    - Master
    - Worker
    - Volumes
    - MaxTries
    - TTLSecondsAfterFinished
    - PriorityClassName
- HorovodJobStatus