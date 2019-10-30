import yaml
import json


class HyKuFe:
    def __init__(self, name, image, cpu, memory, gpu, replica):
        self.data = {'apiVersion': 'hykufe.com/v1alpha1', 'kind': 'HorovodJob', 'metadata': {'name': name, 'labels': {'volcano.sh/job-type': 'Horovod'}}, 'spec': {'schedulerName': 'volcano', 'dataShareMode': {'nfsMode': {'ipAddress': '10.233.96.94', 'path': '/volume'}}, 'dataSources': [{'name': 's3-secret', 's3Source': {'s3SecretName': 's3-secret', 'directory': 'data'}}], 'volumes': [{'volumeClaimName': 'data-volume', 'mountPath': '/data', 'volumeClaim': {'accessModes': ['ReadWriteMany'], 'storageClassName': 'manual', 'resources': {'requests': {'storage': '20Gi'}}, 'volumeMode': 'FileSystem'}}], 'master': {'replicas': 1, 'name': 'master', 'template': {'spec': {'containers': [{'command': ['/bin/bash', '-c', 'set -o pipefail;\nWORKER_HOST=`cat /etc/volcano/worker.host | tr "\\n" ","`;\nmkdir -p /var/run/sshd; /usr/sbin/sshd;\nmkdir -p /result/log;\nsleep 10;\nmpiexec --allow-run-as-root --host ${WORKER_HOST} -np 2 python /examples/tensorflow2_mnist.py 2>&1 | tee /result/log/mpi_log;\n'], 'image': image, 'name': 'master', 'ports': [{'containerPort': 22, 'name': 'job-port'}], 'resources': {'requests': {'cpu': cpu, 'memory': memory, 'nvidia.com/gpu': gpu}, 'limits': {'cpu': cpu, 'memory': memory, 'nvidia.com/gpu': gpu}}}], 'restartPolicy': 'OnFailure', 'imagePullSecrets': [{'name': 'default-secret'}]}}}, 'worker': {'replicas': replica, 'name': 'worker', 'template': {'spec': {'containers': [{'command': ['/bin/sh', '-c', 'mkdir -p /var/run/sshd; /usr/sbin/sshd -D;\n'], 'image': image, 'name': 'worker', 'ports': [{'containerPort': 22, 'name': 'job-port'}], 'resources': {'requests': {'cpu': cpu, 'memory': memory, 'nvidia.com/gpu': gpu}, 'limits': {'cpu': cpu, 'memory': memory, 'nvidia.com/gpu': gpu}}}], 'restartPolicy': 'OnFailure', 'imagePullSecrets': [{'name': 'default-secret'}]}}}}}

    def __str__(self):
        return json.dumps(self.data)


class HyKuFeBuilder:
    def __init__(self):
        self.name = "horovod-job-example"
        self.image = "horovod/horovod:0.18.2-tf2.0.0-torch1.3.0-mxnet1.5.0-py3.6-gpu"
        self.cpu = "2000m"
        self.memory = "4096Mi"
        self.gpu = 1
        self.replica = 2

    def setName(self, name):
        self.name = name
        return self

    def setImage(self, image):
        self.image = image
        return self

    def setCpu(self, cpu):
        self.cpu = cpu
        return self

    def setMemory(self, memory):
        self.memory = memory
        return self

    def setGpu(self, gpu):
        self.gpu = gpu
        return self

    def setReplica(self, replica):
        self.replica = replica
        return self

    def build(self):
        return HyKuFe(self.name, self.image, self.cpu, self.memory, self.gpu, self.replica)


# def readFunc():
#     yaml.dump(yaml.load(open('template.yaml'), Loader=yaml.FullLoader), open('result.yaml', 'w'))
