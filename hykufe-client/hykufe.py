from kubernetes import client, config, utils
import boto3
import yaml
import json
from pprint import pprint

class HyKuFe:
    def __init__(self, accessKey, secretKey, s3BucketName, s3Directory, name, image, cpu, memory, gpu, replica):
        
        config.load_kube_config()
        configuration = client.Configuration()
        secrets = client.CoreV1Api().list_namespaced_secret("default").items
        for item in secrets:
            if "hykufe" in item.metadata.name:
                configuration.api_key['authorization'] = item.data['token']
                break
        configuration.api_key_prefix['authorization'] = 'Bearer'
                
        configuration.verify_ssl = False
        configuration.host = "https://172.16.100.100:6443"

        self.api_instance = client.CustomObjectsApi(client.ApiClient(configuration))


        # about yaml
        self.data = yaml.load(open('template.yaml'), Loader=yaml.FullLoader)

        self.data['spec']['dataSources'][0]['s3Source']['name'] = s3BucketName
        self.data['spec']['dataSources'][0]['s3Source']['directory'] = s3Directory

        self.data['metadata']['name'] = name
        self.data['spec']['master']['template']['spec']['containers'][0]['image'] \
            = self.data['spec']['worker']['template']['spec']['containers'][0]['image'] \
                = image

        master = self.data['spec']['master']['template']['spec']['containers'][0]['resources']
        masterRequest = master['requests']
        masterLimits = master['limits']

        worker = self.data['spec']['worker']['template']['spec']['containers'][0]['resources']
        workerRequest = worker['requests']
        workerLimits = worker['limits']

        masterRequest['cpu'] = masterLimits['cpu'] = workerRequest['cpu'] = workerLimits['cpu'] = cpu
        masterRequest['memory'] = masterLimits['memory'] = workerRequest['memory'] = workerLimits['memory'] = memory
        masterRequest['gpu'] = masterLimits['gpu'] = workerRequest['gpu'] = workerLimits['gpu'] = gpu
        
        self.data['spec']['worker']['replicas'] = replica
        

        # about aws settings
        self.s3Client = boto3.client('s3', aws_access_key_id=accessKey, aws_secret_access_key=secretKey)

    def __str__(self):
        return json.dumps(self.data)

    def writeYamlFile(self, file_name):
        yaml.dump(self.data, open(file_name, 'w'))

    def createJOB(self):
        group = 'hykufe.com' # str | The custom resource's group name
        version = 'v1alpha1' # str | The custom resource's versionc
        plural = 'horovodjobs' # str | The custom resource's plural name. For TPRs this would be lowercase plural kind.
        namespace = 'default'
        name = ''
        pretty = 'true' # str | If 'true', then the output is pretty printed. (optional)
        api_response = self.api_instance.create_namespaced_custom_object(group, version, namespace, plural, self.data, pretty=pretty)
        
        pprint(api_response)

    def uploadFileToS3(self, filePath):
        self.s3Client.upload_file(filePath, \
            self.data['spec']['dataSources'][0]['s3Source']['name'], \
                self.data['spec']['dataSources'][0]['s3Source']['directory']+'/'+filePath.split('/')[-1])

class HyKuFeBuilder:
    def __init__(self):
        self.s3BucketName = 'storage'
        self.s3Directory = 'data'
        self.name = "horovod-job-example"
        self.image = "horovod/horovod:0.18.2-tf2.0.0-torch1.3.0-mxnet1.5.0-py3.6-gpu"
        self.cpu = "2000m"
        self.memory = "4096Mi"
        self.gpu = 1
        self.replica = 2

    def setS3BucketName(self, s3BucketName):
        self.s3BucketName = s3BucketName
        return self

    def setS3Directory(self, s3Directory):
        self.s3Directory = s3Directory
        return self
        
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

    def build(self, accessKey, secretKey):
        return HyKuFe(accessKey, secretKey, self.s3BucketName, self.s3Directory, self.name, self.image, self.cpu, self.memory, self.gpu, self.replica)


# def readFunc():
#     yaml.dump(yaml.load(open('template.yaml'), Loader=yaml.FullLoader), open('result.yaml', 'w'))
