from pprint import pprint

import yaml
from kubernetes import client,watch
from kubernetes.client.rest import ApiException
from object_definition import *

import requests

class TrainingManager:
    # 사용자 이름
    username = None

    global_name = None

    # replicas 개수
    replicas = None

    # namespace
    namespace = None

    # image
    image = None

    #
    entry_point = None
    train_data_dir = None
    test_data_dir = None
    s3_url = "http://ywj-horovod.s3.ap-northeast-2.amazonaws.com"


    # 서버 정보
    token = None
    kube_master_ip = None


    api_client = None

    core_v1_api_instance = None
    apps_v1_api_instance = None
    batch_v1_api_instance = None

    # Definition
    worker_service_def = None
    master_service_def = None
    statefulset_def = None
    job_def = None
    configmap_def = None
    secret_def = None



    def __init__(self, username, token, kube_master_ip, entry_point, train_data_dir=None, test_data_dir=None, replicas = 2, namespace = "default", image="uber/horovod:0.12.1-tf1.8.0-py3.5"):

        self.entry_point = entry_point

        if token is None or kube_master_ip is None:
            raise Exception()

        self.token = token
        self.kube_master_ip = kube_master_ip


        if username is not None:
            self.username = username
            self.global_name = username + "-horovod"

        if replicas is not None:
            self.replicas = replicas

        if namespace is not None:
            self.namespace = namespace

        if image is not None:
            self.image = image

        # 마스터 서버에 연결합니다.
        configuration = client.Configuration()
        configuration.host = self.kube_master_ip

        # 일단 테스트를 위해서 insecure 옵션 추가
        configuration.verify_ssl = False

        configuration.api_key = {"authorization": "Bearer " + self.token }

        # ApiClient를 생성합니다.
        self.api_client = client.ApiClient(configuration)

        self.core_v1_api_instance = client.CoreV1Api(self.api_client)
        self.apps_v1_api_instance = client.AppsV1Api(self.api_client)
        self.batch_v1_api_instance = client.BatchV1Api(self.api_client)

    def sendDataToS3(self):
        if self.entry_point is None:
            raise Exception("Entry Point를 입력하세요.")

        # Todo train, test data 전송 짜기
        # TOdo: 일단 S3를 Pulbic Access로 해놨는데 이거 수정하기 ... 제발

        # entry_point 전송
        entry_point_file = open(self.entry_point, 'rb')

        upload = {'file': entry_point_file}

        res = requests.put(self.s3_url + "/horovod/" + self.username + "/" + "train.py", data=entry_point_file)

        if res.status_code != 200:
            raise Exception("파일 전송 실패")

        print(res.status_code)

    def createDefinition(self):

        self.statefulset_def = createStatefulSet(self.username, self.replicas, self.image)
        self.job_def = createJob(self.username, self.image, self.replicas)
        self.configmap_def = createConfigmap(self.username, self.replicas)
        self.secret_def = createSecret(self.username)
        self.worker_service_def = createService(self.username, "master")
        self.master_service_def = createService(self.username, "worker")



    def createAllObject(self):
        try:
            self.core_v1_api_instance = client.CoreV1Api(self.api_client)
            self.apps_v1_api_instance = client.AppsV1Api(self.api_client)
            self.batch_v1_api_instance = client.BatchV1Api(self.api_client)

            # configmap 생성
            api_response = self.core_v1_api_instance.create_namespaced_config_map(self.namespace, self.configmap_def)
            pprint(api_response)

            # Secret 생성
            api_response = self.core_v1_api_instance.create_namespaced_secret(self.namespace, self.secret_def)
            pprint(api_response)

            # StatefulSet 생성
            api_response = self.apps_v1_api_instance.create_namespaced_stateful_set(self.namespace, self.statefulset_def)
            pprint(api_response)



            # Job 생성
            api_response = self.batch_v1_api_instance.create_namespaced_job(self.namespace, self.job_def)
            pprint(api_response)

            # Master Service 생성
            api_response = self.core_v1_api_instance.create_namespaced_service(self.namespace, self.master_service_def)
            pprint(api_response)

            # Worker Service 생성
            api_response = self.core_v1_api_instance.create_namespaced_service(self.namespace, self.worker_service_def)
            pprint(api_response)

        except ApiException as e:
            print(e)
            raise e

    def deleteAllObject(self):
        try:
            #master_pod_name = self.core_v1_api_instance.list_namespaced_pod(namespace=self.namespace, label_selector='job-name={}'.format(self.global_name)).items[0].metadata.name

            self.core_v1_api_instance.delete_namespaced_config_map(self.global_name, self.namespace)
            self.core_v1_api_instance.delete_namespaced_secret(self.global_name, self.namespace)
            self.apps_v1_api_instance.delete_namespaced_stateful_set(self.global_name, self.namespace)
            self.core_v1_api_instance.delete_namespaced_service(self.global_name + "-worker", self.namespace)
            self.core_v1_api_instance.delete_namespaced_service(self.global_name + "-master", self.namespace)
            self.batch_v1_api_instance.delete_namespaced_job(self.global_name, self.namespace)

            #self.core_v1_api_instance.delete_namespaced_pod(name=master_pod_name, namespace=self.namespace)
        except ApiException as e:
            pass

    def runTrain(self):
        # 모든 팟과 설정파일을 올리고
        try:
            self.sendDataToS3()
            self.createDefinition()
            self.createAllObject()
            # 로그 출력

            try:
                # Master 파드의 이름을 찾는다.
                master_pod_name = self.core_v1_api_instance.list_namespaced_pod(namespace = self.namespace, label_selector='job-name={}'.format(self.global_name)).items[0].metadata.name
                print(master_pod_name)
                stream = watch.Watch().stream(self.core_v1_api_instance.read_namespaced_pod_log, name=master_pod_name , namespace=self.namespace)

                for event in stream:
                    print(event)
            except ApiException as e:
                print('Found exception in reading the logs')

        except Exception as e:
            print(e)
        finally:
            pass
            #self.deleteAll()
        # 로그와 상태를 보고 끝났는지 확인

        #

train_script = "tensorflow2_mnist.py"

tm = TrainingManager("ywj", token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6Inl3ai1zYS10b2tlbi1sbnh3bSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJ5d2otc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI1YzRjMGYzZi1iOWYwLTExZTktYTNhYy1mYTE2M2UwOTllM2YiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDp5d2otc2EifQ.dxalsvH1kTLbD6FULnp3K6ukxebqpkxAM3myq2dYJIETKL1VGfGOPxVmY5jEjiXE_Lgb7uKHu7kckXDoVieeOljGxUPE1wGRNlicWHm2BKu57QQ8wQimaERJTrxHbxuh6d-U90lU28Yg4Y-7BmGunp2VJeBAZ2ajsNRMKw-u1l38glaOQlKFtNby94KzcaNA0jPsTXOFKKf7UMmJiUXhlJ1Pf-RhqHZ72jV4ZXr1OT9hpFJBEqUP9C4iC5k8prc36IyUr6-9mDwWlFV6VhPvtGt9DkoAF7DVYvj0MrZYrreVtFfZTuv5LMdgce6hjWI3e0wLkMnBL8o0K57Z1PC3UQ",
                     replicas=2,
                     kube_master_ip="203.254.143.253:8080",
                     entry_point=train_script,
                     image="horovod/horovod:0.16.4-tf1.14.0-torch1.1.0-mxnet1.4.1-py3.6"
                     )

tm.sendDataToS3()
tm.deleteAllObject()
tm.runTrain()
