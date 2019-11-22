# HyKuFe
# Hykufe Operator
hykufe.com/hykufe-operator

## Prerequisite
### NFS Server
만약 분산학습을 위한 NFS가 없는 경우 임시 nfs server를 설정할 수 있다.
```
$ kubectl create -f nfs-server/nfs-server-pv.yml
$ kubectl create -f nfs-server/nfs-server-pvc.yml
$ kubectl create -f nfs-server/nfs-server-deployment.yml
```

AWS 인스턴스를 동적으로 프로비저닝하고 하이브리드 클러스터를 구성하려면 AWS 액세스 키가 필요하다.



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


## Skaffold를 이용한 자동 배포 및 테스트
### prerequisite
개발환경에 배포할 클러스터에 접근가능해야 합니다.
kubernetes api server에 직접 접근하려면 원격 클러스터의 접근 정보와 private key를 직접 개발환경에 복사해서 사용하는 방법이 있습니다.(내가 쓰고 있음)
이 경우 이 키는 개발환경에 대한 신뢰할 수 있는 사용자로 등록되어있지 않으므로 kubeadm을 이용하여 개발환경에 대해서 신뢰할 수 있는 사용자로 설정하여 다시 키를 발급받는 방법과 --insecure-skip-tls를 사용하여 인증을 건너뛰는 방법이 있습니다.


### skaffold 설치
- Macos - brew install skaffold

### skaffold 실행
hykufe-operator 디렉토리아래 skaffold.yaml파일을 열어서 image부분을 적절한 도커 레지스트리로 설정합니다.
설정이 완료되면 다음 명령어를 사용하면 자동으로 build 스크립트를 호출하고 deploy 아래 있는 yaml파일들을 배포할 수 있습니다.

``` shell
$ skaffold dev
```

소스코드의 변경을 감지하므로 만약 소스코드의 변경이 있을 경우 자동으로 다시 빌드를 하고 배포를 진행합니다.
배포가 완료되면 skaffold는 operator의 pod를 식별하여 자동으로 로그를 출력해주게 됩니다. 

### Issue
*_crd.yaml 파일은 버그로 인해 자동으로 배포가 되지 않도록 설정했습니다. 따라서 type에 변화가 있다면 반드시 수동으로 지웠다가 다시 생성해주어야 합니다.

## Custrom Resource Definition(CRD)
### HorovodJob
- name
- type
- HorovodJobSpec
    - Master
    - Worker
    - DataShareMode
        - NFSMode
            - IPAddress
            - Directory
            - AccessMode
    - DataSource[]
        - DataSourceSpec
    - Volumes
    - MaxTries
    - TTLSecondsAfterFinished
    - PriorityClassName
- HorovodJobStatus