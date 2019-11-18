from kubernetes import client, config, utils
import yaml
import json
from pprint import pprint

class HyKuFe:
    def __init__(self, name, image, cpu, memory, gpu, replica):
        
        # about kubernetes settings
        configuration = client.Configuration()
        # configuration.api_key['authorization'] = 'eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tZHY1N3ciLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImZmNWI3MGMzLWU3NTMtMTFlOS05NjI0LTcwODVjMjAyYmUyOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.JsJj2cp1kWyeZLz4Tm0NiCH7hwQOvlf1PtTXWX1k0drjev1LmXMJOIQk6GSAhlCK-eRUa2rLVENLtC6Tlo_hXVfl7frHDL1N6jjb3ZBpR4hvcxkCXPvkkr2mjIxGCKXcsPhGiGjZ1DazFxttT6Vh9DdZ04Oa8TiDP76Dqjo5Pfv3VvdV1YPLN8WXYEN-IJE7Et-tYgEz5eepxXACjISR6VsFly0os9F6RMLnkfxZxP-JOpZspmQPlnfTJXtpLRZGiLsAC3A7tEp2SLnHtPpmveIixK47HIpQXWNsTwOUZG9oTfDjRXODFAjiIn9dMRREfT1qjK4Wl6ovjyPGcxW0cA'
        # # # Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
        # configuration.api_key_prefix['authorization'] = 'Bearer'
        
        # configuration.api_key = {"authorization": "Bearer " + "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tZHY1N3ciLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImZmNWI3MGMzLWU3NTMtMTFlOS05NjI0LTcwODVjMjAyYmUyOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.JsJj2cp1kWyeZLz4Tm0NiCH7hwQOvlf1PtTXWX1k0drjev1LmXMJOIQk6GSAhlCK-eRUa2rLVENLtC6Tlo_hXVfl7frHDL1N6jjb3ZBpR4hvcxkCXPvkkr2mjIxGCKXcsPhGiGjZ1DazFxttT6Vh9DdZ04Oa8TiDP76Dqjo5Pfv3VvdV1YPLN8WXYEN-IJE7Et-tYgEz5eepxXACjISR6VsFly0os9F6RMLnkfxZxP-JOpZspmQPlnfTJXtpLRZGiLsAC3A7tEp2SLnHtPpmveIixK47HIpQXWNsTwOUZG9oTfDjRXODFAjiIn9dMRREfT1qjK4Wl6ovjyPGcxW0cA" }
        
        # configuration.verify_ssl = False
        # configuration.host = "https://172.16.100.100:6443"
        configuration.host = "127.0.0.1:8001"

        self.api_instance = client.CustomObjectsApi(client.ApiClient(configuration))
        
        # about yaml
        self.data = yaml.load(open('template.yaml'), Loader=yaml.FullLoader)
        
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
