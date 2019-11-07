import yaml
import json


class HyKuFe:
    def __init__(self, name, image, cpu, memory, gpu, replica):
        
        self.data = yaml.load(open('template.yaml'), Loader=yaml.FullLoader)
        
        self.data['metadata']['name'] = name
        self.data['spec']['master']['template']['spec']['containers'][0]['image'] \
            = self.data['spec']['worker']['template']['spec']['containers'][0]['image'] \
                = image

        self.data['spec']['master']['template']['spec']['containers'][0]['resources']['requests']['cpu'] \
            = self.data['spec']['master']['template']['spec']['containers'][0]['resources']['limits']['cpu'] \
                = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']['requests']['cpu'] \
                    = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']['limits']['cpu'] \
                        = cpu

        self.data['spec']['master']['template']['spec']['containers'][0]['resources']['requests']['memory'] \
            = self.data['spec']['master']['template']['spec']['containers'][0]['resources']['limits']['memory'] \
                = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']['requests']['memory'] \
                    = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']['limits']['memory'] \
                        = memory

        self.data['spec']['master']['template']['spec']['containers'][0]['resources']['requests']['nvidia.com/gpu'] \
            = self.data['spec']['master']['template']['spec']['containers'][0]['resources']['limits']['nvidia.com/gpu'] \
                = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']['requests']['nvidia.com/gpu'] \
                    = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']['limits']['nvidia.com/gpu'] \
                        = gpu
        
        self.data['spec']['worker']['replicas'] = replica

    def __str__(self):
        return json.dumps(self.data)

    def writeYamlFile(self, file_name):
        yaml.dump(self.data, open(file_name, 'w'))


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

    def setCPU(self, cpu):
        self.cpu = cpu
        return self

    def setMemory(self, memory):
        self.memory = memory
        return self

    def setGPU(self, gpu):
        self.gpu = gpu
        return self

    def setReplica(self, replica):
        self.replica = replica
        return self

    def build(self):
        return HyKuFe(self.name, self.image, self.cpu, self.memory, self.gpu, self.replica)


# def readFunc():
#     yaml.dump(yaml.load(open('template.yaml'), Loader=yaml.FullLoader), open('result.yaml', 'w'))
