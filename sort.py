from pprint import pprint

import yaml
from kubernetes import client,watch
from kubernetes.client.rest import ApiException
from object_definition import *
from time import sleep
import requests

class TrainingClusterManager:
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


    # Host Network
    is_host_network = None
    ssh_port = None

    # 서버 정보
    token = None
    kube_master_ip = None


    # Kubernetes Client API
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



    def __init__(self, username, token, kube_master_ip, entry_point, train_data_dir=None, test_data_dir=None, replicas = 2, namespace = "default", image="uber/horovod:0.12.1-tf1.8.0-py3.5", is_host_network=False, ssh_port="22"):

        self.global_name = "{username}-horovod".format(username = username)
        self.is_host_network = is_host_network
        self.ssh_port = ssh_port
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

        self.statefulset_def = createStatefulSet(self.username, self.replicas, self.image, is_host_network=self.is_host_network, ssh_port=self.ssh_port)
        self.job_def = createJob(self.username, self.image, self.replicas, is_host_network=self.is_host_network, ssh_port=self.ssh_port)
        self.configmap_def = createConfigmap(self.username, self.replicas)
        self.secret_def = createSecret(self.username)
        self.worker_service_def = createService(self.username, "master", ssh_port=self.ssh_port)
        self.master_service_def = createService(self.username, "worker", ssh_port=self.ssh_port)



    def createAllObject(self):
        try:
            self.core_v1_api_instance = client.CoreV1Api(self.api_client)
            self.apps_v1_api_instance = client.AppsV1Api(self.api_client)
            self.batch_v1_api_instance = client.BatchV1Api(self.api_client)

            # configmap 생성
            api_response = self.core_v1_api_instance.create_namespaced_config_map(self.namespace, self.configmap_def)
            print("created config map")

            # Secret 생성
            api_response = self.core_v1_api_instance.create_namespaced_secret(self.namespace, self.secret_def)
            print("created secret")

            # StatefulSet 생성
            api_response = self.apps_v1_api_instance.create_namespaced_stateful_set(self.namespace, self.statefulset_def)
            print("created statefulset")



            # Job 생성
            api_response = self.batch_v1_api_instance.create_namespaced_job(self.namespace, self.job_def)
            print("created job")

            # Master Service 생성
            api_response = self.core_v1_api_instance.create_namespaced_service(self.namespace, self.master_service_def)
            print("created master service")

            # Worker Service 생성
            api_response = self.core_v1_api_instance.create_namespaced_service(self.namespace, self.worker_service_def)
            print("created worker service")

        except ApiException as e:
            print(e)
            raise e

    def deleteAllObject(self):

        # Configmap을 제거합니다.
        try:
            self.core_v1_api_instance.delete_namespaced_config_map(self.global_name, self.namespace)
        except ApiException as e:
            print(e)

        # Secret을 제거합니다.
        try:
            self.core_v1_api_instance.delete_namespaced_secret(self.global_name, self.namespace)
        except ApiException as e:
            print(e)

        # Job을 제거합니다.
        try:
            self.batch_v1_api_instance.delete_namespaced_job(self.global_name, self.namespace)
        except ApiException as e:
            print(e)

        master_pod_name_list = self.core_v1_api_instance.list_namespaced_pod(namespace=self.namespace, label_selector='job-name={0}'.format(self.global_name))

        # Job Pod를 제거합니다.
        for master_pod in master_pod_name_list.items:
            master_pod_name = master_pod.metadata.name
            self.core_v1_api_instance.delete_namespaced_pod(name=master_pod_name, namespace=self.namespace)

        # Statefulset을 제거합니다.
        try:
            self.apps_v1_api_instance.delete_namespaced_stateful_set(self.global_name, self.namespace)
        except ApiException as e:
            print(e)

        # worker Service를 제거합니다.
        try:
            self.core_v1_api_instance.delete_namespaced_service(self.global_name + "-worker", self.namespace)
        except ApiException as e:
            print(e)

        # Master Service를 제거합니다.
        try:
            self.core_v1_api_instance.delete_namespaced_service(self.global_name + "-master", self.namespace)
        except ApiException as e:
            print(e)




    def runTrain(self):
        # 모든 팟과 설정파일을 올리고
        try:
            self.sendDataToS3()
            self.createDefinition()
            self.createAllObject()
            # 로그 출력


            sleep(1)

            # Master 파드의 이름을 찾는다.
            master_pod_name = self.core_v1_api_instance.list_namespaced_pod(namespace = self.namespace, label_selector='job-name={}'.format(self.global_name)).items[0].metadata.name

            # Master 파드의 활성 여부를 3초 간격으로 10번 확인한다.
            master_pod = None
            for i in range(0, 500):
                sleep(3)

                master_pod = self.core_v1_api_instance.read_namespaced_pod_status(name=master_pod_name, namespace=self.namespace)

                print("Waiting to Ready")
                if master_pod.status.phase == "Running":
                    break
            if master_pod.status.phase != "Running":
                print("Pod 생성에 실패했습니다.")

            master_pod_name = master_pod.metadata.name

            while True:
                master_job_status = self.batch_v1_api_instance.read_namespaced_job_status(name=self.global_name, namespace=self.namespace)
                #print(master_job_status.status.failed, master_job_status.status.succeeded)


                # Job이 끝났을 때
                if master_job_status.status.active != 1:
                    if master_job_status.status.succeeded == 1:
                        print("학습이 완료되었습니다.")
                    if master_job_status.status.failed == 1:
                        print("학습이 실패했습니다. 로그를 확인해주세요.")
                    break

                # Master Pod의 로그를 출력한다.
                stream = watch.Watch().stream(self.core_v1_api_instance.read_namespaced_pod_log, name=master_pod_name, namespace=self.namespace)
                for event in stream:
                    print(event)


        except Exception as e:
            print(e)
        finally:
            pass
            #self.deleteAllObject()


