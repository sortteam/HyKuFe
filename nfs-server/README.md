# NFS-Server
분산학습을 위한 임시 NFS-server입니다.
```shell script
# create directory for pv(Hostpath) (ubuntu)
$ mkdir /home/ubuntu/pv

$ kubectl create -f nfs-server-pv.yml
$ kubectl create -f nfs-server-pvc.yml
$ kubectl create -f nfs-server-deployment.yml


# get pod ClusterIP 
$ kubectl describe pod nfs-server-?????

```

이제 Container Spec에 nfs volume을 Cluster IP로 선언하고 마운트해서 쓰면 됩니다.


## ISSUE
- Service를 생성하고 nfs-server 파드에 endpoint까지 활성화되었지만 해당 Service ClusterIP로 nfs mount 시도시 에러 발생
