import yaml
from kubernetes import client

def createConfigmap(username, replicas):
    configmap_name = username + "-horovod"
    hostfile = configmap_name + "-master slots=1\n"
    for i in range(0, replicas):
        hostfile += configmap_name + "-" + str(i) + "." + configmap_name + "-worker slots=1\n"

    configmap = client.V1ConfigMap()
    configmap.apiVersion = "V1"
    configmap.kind = "ConfigMap"
    configmap.metadata = client.V1ObjectMeta(name=configmap_name,
                                             labels={
                                                 "user": username,
                                                 "app": 'horovod'
                                             })
    # 외부 쉘파일을 불러온다.
    # 좋은 방법 있음?
    with open("sh.yaml", 'r') as f:
        try:
            data = yaml.safe_load(f)
            configmap.data = {
                "hostfile.config": hostfile,
                "ssh.readiness": data['data']["ssh.readiness"],
                "master.run": data['data']['master.run'],
                "master.waitWorkerReady": data['data']['master.waitWorkerReady'],
                "worker.run": data['data']['worker.run']
            }
        except yaml.YAMLError as e:
            raise e

    return configmap


def createSecret(username):
    secret_name = username + "-horovod"

    secret = client.V1Secret()
    # secret.kind = 'secret'
    secret.metadata = client.V1ObjectMeta(name=secret_name, labels={
        "app": 'horovod',
        "user": username
    })

    secret.type = "Obaque"
    secret.data = {
#         "host-key": "-----BEGIN RSA PRIVATE KEY-----\
# MIIEowIBAAKCAQEAphzfpjDIEF2CqlDbLhF+t82GhJKSb8twTu0HJY9adfGBSDDP\
# wXWUUZYTDTHY4mjXgApYj443mC6P2ULr+xE/lQAP6xG3at9r2R6Adx367LrnWtp1\
# 7Ro641bdtg9wkY2fDnBS+aAHiYxUEpe03K6PsGTh4TrASZs3aZHqRAxU2WYH27aG\
# hx3HQ1KRGGLm5pJJnoxM2ZFac/RkuGH5Lz5JACbgYRHSV8ybfy77yy7o2KYVgUlW\
# eeuOaTquRHmYJvuylYGzKeAi7k+vpKBP2Fo/CY7RteFQKgI5BnlhGCmbR9WTuoL+\
# zb9ngdHRnQjhdiMyaUH5DBXMBZwQjzyodTz6jQIDAQABAoIBACcavPuemDpiCRSX\
# HEHoFHCojXZAGwD+X131Jq2M5brGM60O8JmWWGgscCe3CFukWrbluJty21uT+oEm\
# 4+6izNkCvryT2x3porXmHE/uWtfH2BbnPsOmXR6PoHnvgIyDTmJTxvTE24Fh65jE\
# 5erdnS3lUdd3wTSSuaS8mO2UCZVzsmVGDUs/dl7hlFkkJ8oXWIImh6yCW/9Ti8uV\
# yC/rQi7Lkctfz0tc8RKNeLXk8D3qQzSpSOWlHh4BAXsuWNoZVMLXP2myqqnaCCy3\
# o7XjUQ56u/hIwAE58RC+usFH+KrPVAiexnd298SVuo/a7X3rViSMbzBVnq8cM89O\
# Rcx324ECgYEA0D3iVv9STQs5yB2d0aL61WVbkmUzHDxihyTzEEcNQxkMO1P6EnG8\
# aSAlx7Z+0c/0Q7q/6xad9fweVRbSOyw52DIhJL+aE11nHx9u1VP/Mc9/SeYn05jy\
# A81VLHWE2FxHPQWccwv2M+yCLT8OwLkYYlEzBZrLHehF3HVSuCWyVlUCgYEAzDWM\
# udRxYfRDvbeol/sxdlKPi3qBz3quJ6mTWi1cdl7gwgHZYZi8KVpY/u08KuYBpbjK\
# G20owy9VCg8PCXmgOJGBkpS3wF/p7LufQIxWKFfukaaguYfOIQBUZowqEnR7tQX2\
# KKALdE3efZQ+zSk3WRbjJjR7FQv3LPDtXyyaG1kCgYAR7JGz3UwvN30kvW/dIIMo\
# pQ3Jvw40KvpsGYEWcJcypFBKNwM6XTHTdqHp28p0ssqandNxH8Q+7RGLT2iPEVJ1\
# SnNR33AapJqAskru78jyd6LEMJxS+UIzk5P2PLNPkDnNhdMej/QEKiJWVKwnaIcx\
# xz12CQncrCZ/QFX7ZbtA9QKBgHfjiHchLl/f1FVxmd2AcA2TcwrkJYn18IEAoa3z\
# q7EjCrlb9I/D59QvYshn50cYOiddUerAL4pII5kANkfNzC7p3jR8c1TR+rgtftWa\
# joqo9Ts1pG7IOFBPrT13VMv47xfcJCS9sXvaq6D2g9hXNlNriHhJn2k/2SHdYL7b\
# pK4hAoGBALTSDuRAn6+heQQfRReMNjXxqADdq0clbwh7XFDEOuCXbLgxqq36zA/Z\
# k+AJ29wuD/LUz0vqjmoGFwJhPh+2qbLiL6mD9F1rQWCCCeGMo3Zwvcw5RIGp6K6C\
# 5qSHjhVmwG9JJlk2XRwxdxyPoGRsVI8lhX92LFu8iKAv7zpDLmsX\
# -----END RSA PRIVATE KEY-----",
        "host-key": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBcGh6ZnBqRElF\
RjJDcWxEYkxoRit0ODJHaEpLU2I4dHdUdTBISlk5YWRmR0JTRERQd1hXVVVaWVREVEhZNG1qWGdB\
cFlqNDQzbUM2UDJVTHIreEUvbFFBUDZ4RzNhdDlyMlI2QWR4MzY3THJuV3RwMTdSbzY0MWJkdGc5\
d2tZMmZEbkJTK2FBSGlZeFVFcGUwM0s2UHNHVGg0VHJBU1pzM2FaSHFSQXhVMldZSDI3YUdoeDNI\
UTFLUkdHTG01cEpKbm94TTJaRmFjL1JrdUdINUx6NUpBQ2JnWVJIU1Y4eWJmeTc3eXk3bzJLWVZn\
VWxXZWV1T2FUcXVSSG1ZSnZ1eWxZR3pLZUFpN2srdnBLQlAyRm8vQ1k3UnRlRlFLZ0k1Qm5saEdD\
bWJSOVdUdW9MK3piOW5nZEhSblFqaGRpTXlhVUg1REJYTUJad1FqenlvZFR6NmpRSURBUUFCQW9J\
QkFDY2F2UHVlbURwaUNSU1hIRUhvRkhDb2pYWkFHd0QrWDEzMUpxMk01YnJHTTYwTzhKbVdXR2dz\
Y0NlM0NGdWtXcmJsdUp0eTIxdVQrb0VtNCs2aXpOa0N2cnlUMngzcG9yWG1IRS91V3RmSDJCYm5Q\
c09tWFI2UG9IbnZnSXlEVG1KVHh2VEUyNEZoNjVqRTVlcmRuUzNsVWRkM3dUU1N1YVM4bU8yVUNa\
VnpzbVZHRFVzL2RsN2hsRmtrSjhvWFdJSW1oNnlDVy85VGk4dVZ5Qy9yUWk3TGtjdGZ6MHRjOFJL\
TmVMWGs4RDNxUXpTcFNPV2xIaDRCQVhzdVdOb1pWTUxYUDJteXFxbmFDQ3kzbzdYalVRNTZ1L2hJ\
d0FFNThSQyt1c0ZIK0tyUFZBaWV4bmQyOThTVnVvL2E3WDNyVmlTTWJ6QlZucThjTTg5T1JjeDMy\
NEVDZ1lFQTBEM2lWdjlTVFFzNXlCMmQwYUw2MVdWYmttVXpIRHhpaHlUekVFY05ReGtNTzFQNkVu\
RzhhU0FseDdaKzBjLzBRN3EvNnhhZDlmd2VWUmJTT3l3NTJESWhKTCthRTExbkh4OXUxVlAvTWM5\
L1NlWW4wNWp5QTgxVkxIV0UyRnhIUFFXY2N3djJNK3lDTFQ4T3dMa1lZbEV6QlpyTEhlaEYzSFZT\
dUNXeVZsVUNnWUVBekRXTXVkUnhZZlJEdmJlb2wvc3hkbEtQaTNxQnozcXVKNm1UV2kxY2RsN2d3\
Z0haWVppOEtWcFkvdTA4S3VZQnBiaktHMjBvd3k5VkNnOFBDWG1nT0pHQmtwUzN3Ri9wN0x1ZlFJ\
eFdLRmZ1a2FhZ3VZZk9JUUJVWm93cUVuUjd0UVgyS0tBTGRFM2VmWlErelNrM1dSYmpKalI3RlF2\
M0xQRHRYeXlhRzFrQ2dZQVI3Skd6M1V3dk4zMGt2Vy9kSUlNb3BRM0p2dzQwS3Zwc0dZRVdjSmN5\
cEZCS053TTZYVEhUZHFIcDI4cDBzc3FhbmROeEg4USs3UkdMVDJpUEVWSjFTbk5SMzNBYXBKcUFz\
a3J1NzhqeWQ2TEVNSnhTK1VJems1UDJQTE5Qa0RuTmhkTWVqL1FFS2lKV1ZLd25hSWN4eHoxMkNR\
bmNyQ1ovUUZYN1pidEE5UUtCZ0hmamlIY2hMbC9mMUZWeG1kMkFjQTJUY3dya0pZbjE4SUVBb2Ez\
enE3RWpDcmxiOUkvRDU5UXZZc2huNTBjWU9pZGRVZXJBTDRwSUk1a0FOa2ZOekM3cDNqUjhjMVRS\
K3JndGZ0V2Fqb3FvOVRzMXBHN0lPRkJQclQxM1ZNdjQ3eGZjSkNTOXNYdmFxNkQyZzloWE5sTnJp\
SGhKbjJrLzJTSGRZTDdicEs0aEFvR0JBTFRTRHVSQW42K2hlUVFmUlJlTU5qWHhxQURkcTBjbGJ3\
aDdYRkRFT3VDWGJMZ3hxcTM2ekEvWmsrQUoyOXd1RC9MVXowdnFqbW9HRndKaFBoKzJxYkxpTDZt\
RDlGMXJRV0NDQ2VHTW8zWnd2Y3c1UklHcDZLNkM1cVNIamhWbXdHOUpKbGsyWFJ3eGR4eVBvR1Jz\
Vkk4bGhYOTJMRnU4aUtBdjd6cERMbXNYCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==",        #"host-key-pub": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCmHN+mMMgQXYKqUNsuEX63zYaEkpJvy3BO7Qclj1p18YFIMM/BdZRRlhMNMdjiaNeACliPjjeYLo/ZQuv7ET+VAA/rEbdq32vZHoB3Hfrsuuda2nXtGjrjVt22D3CRjZ8OcFL5oAeJjFQSl7Tcro+wZOHhOsBJmzdpkepEDFTZZgfbtoaHHcdDUpEYYubmkkmejEzZkVpz9GS4YfkvPkkAJuBhEdJXzJt/LvvLLujYphWBSVZ5645pOq5EeZgm+7KVgbMp4CLuT6+koE/YWj8JjtG14VAqAjkGeWEYKZtH1ZO6gv7Nv2eB0dGdCOF2IzJpQfkMFcwFnBCPPKh1PPqN ubuntu@node1"
        "host-key-pub": "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFDbUhOK21NTWdRWFlLcVVOc3VF\
WDYzellhRWtwSnZ5M0JPN1FjbGoxcDE4WUZJTU0vQmRaUlJsaE1OTWRqaWFOZUFDbGlQamplWUxv\
L1pRdXY3RVQrVkFBL3JFYmRxMzJ2WkhvQjNIZnJzdXVkYTJuWHRHanJqVnQyMkQzQ1JqWjhPY0ZM\
NW9BZUpqRlFTbDdUY3JvK3daT0hoT3NCSm16ZHBrZXBFREZUWlpnZmJ0b2FISGNkRFVwRVlZdWJt\
a2ttZWpFelprVnB6OUdTNFlma3ZQa2tBSnVCaEVkSlh6SnQvTHZ2TEx1allwaFdCU1ZaNTY0NXBP\
cTVFZVpnbSs3S1ZnYk1wNENMdVQ2K2tvRS9ZV2o4Smp0RzE0VkFxQWprR2VXRVlLWnRIMVpPNmd2\
N052MmVCMGRHZENPRjJJekpwUWZrTUZjd0ZuQkNQUEtoMVBQcU4gdWJ1bnR1QG5vZGUxCg=="
    }
    return secret


def createStatefulSet(username, replicas, image):
    statefulset_name = username + "-horovod"

    statefulset = client.V1StatefulSet()
    #statefulset.api_version="apps/v1beta2"
    statefulset.metadata = client.V1ObjectMeta(name=statefulset_name, labels={
        "app": "horovod",
        "user": username,
        "role": "worker"
    })

    label_selector = client.V1LabelSelector(match_labels={
        "app": "horovod",
        "user": username,
        "role": "worker"
    })

    # Pod template 정의

    pod_template = client.V1PodTemplateSpec()
    pod_template.metadata = client.V1ObjectMeta(labels={
        "app": "horovod",
        "user": username,
        "role": "worker"
    })

    container = client.V1Container(name="worker")
    container.image = image
    container.image_pull_policy = "IfNotPresent"
    container.env = [
        client.V1EnvVar(name="SSHPORT", value="22"),
        client.V1EnvVar(name="USESECRETS", value="true"),
        # TODO: 바꾸기
        client.V1EnvVar(name="ENTRY_POINT", value="train.py")
    ]
    container.ports = [
        client.V1ContainerPort(container_port=22)
    ]
    container.volume_mounts = [
        client.V1VolumeMount(name=statefulset_name + "-cm", mount_path="/horovod/generated"),
        client.V1VolumeMount(name=statefulset_name + "-secret", mount_path="/etc/secret-volume", read_only=True),
        client.V1VolumeMount(name=statefulset_name + "-data", mount_path="/horovod/data")
    ]
    container.command = [
        "/horovod/generated/run.sh"
    ]
    container.readiness_probe = client.V1Probe(_exec=client.V1ExecAction(command=["/horovod/generated/check.sh"]),
                                               initial_delay_seconds=1,
                                               period_seconds=2)

    pod_spec = client.V1PodSpec(containers=[container])
    pod_spec.volumes = [
        client.V1Volume(name=statefulset_name + "-cm",
                        config_map=client.V1ConfigMapVolumeSource(name=statefulset_name,
                                                                  items=[
                                                                      client.V1KeyToPath(key="hostfile.config",
                                                                                         path="hostfile",
                                                                                         mode=438),
                                                                      client.V1KeyToPath(key="ssh.readiness",
                                                                                         path="check.sh",
                                                                                         mode=365),
                                                                      client.V1KeyToPath(key="worker.run",
                                                                                         path="run.sh",
                                                                                         mode=365)

                                                                  ]
                                                                  )),
        client.V1Volume(name=statefulset_name + "-secret",
                        secret=client.V1SecretVolumeSource(secret_name=statefulset_name,
                                                           default_mode=448,
                                                           items=[
                                                               client.V1KeyToPath(key="host-key",
                                                                                  path="id_rsa"),
                                                               client.V1KeyToPath(key="host-key-pub",
                                                                                  path="authorized_keys")
                                                           ]

                                                           )),
        client.V1Volume(name=statefulset_name + "-data",
                        empty_dir=client.V1EmptyDirVolumeSource())

    ]
    pod_spec.subdomain = statefulset_name
    pod_spec.hostname = statefulset_name
    pod_spec.init_containers = [
        client.V1Container(name="download-data",
                           image=image,
                           image_pull_policy="IfNotPresent",
                           command=[
                               "/bin/bash",
                               "-c"
                           ],
                           args=[
                               "curl http://ywj-horovod.s3.ap-northeast-2.amazonaws.com/horovod/" + username + "/train.py > /horovod/data/train.py"
                           ],
                           volume_mounts=[
                               client.V1VolumeMount(name=statefulset_name + "-data", mount_path="/horovod/data")
                           ])
    ]
    pod_template.spec = pod_spec

    statefulset.spec = client.V1StatefulSetSpec(selector=label_selector,
                                                service_name=statefulset_name + "-worker", # https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#pod-identity
                                                pod_management_policy="Parallel",
                                                replicas=replicas,
                                                template=pod_template)
    return statefulset




def createJob(username, image, replicas, train_mode="cpu"):
    job_name = username + "-horovod"

    job = client.V1Job()
    job.metadata = client.V1ObjectMeta(name=job_name,
                                       labels={
                                           "app": "horovod",
                                           "user": username,
                                           "role": "master"
                                       })
    # Job Spec 정의
    job_spec = client.V1JobSpec(template="")
    # job_spec.metadata = client.V1ObjectMeta(name=job_name, labels={
    #     "app": "horovod",
    #     "user": username,
    #     "role": "master"
    # })
    # Job Spec의 Pod Template 정의
    pod_template_spec = client.V1PodTemplateSpec()
    pod_template_spec.restart_policy = "OnFailure"
    pod_template_spec.metadata = client.V1ObjectMeta(name=job_name, labels={
        "app": "horovod",
        "user": username,
        "role": "master"
    })
    pod_template_spec.spec = ""

    # Pod Spec 정의
    pod_spec = client.V1PodSpec(containers=[""], restart_policy="Never")
    # Container 정의
    container = client.V1Container(name=job_name + "-master")

    container.image = image
    container.image_pull_policy = "IfNotPresent"
    container.env = [
        client.V1EnvVar(name="SSHPORT", value="22"),
        client.V1EnvVar(name="USESECRETS", value="true"),
        # TODO: 바꾸기
        client.V1EnvVar(name="ENTRY_POINT", value="train.py")
    ]
    container.ports = [
        client.V1ContainerPort(container_port=22)
    ]
    container.volume_mounts = [
        client.V1VolumeMount(name=job_name + "-cm", mount_path="/horovod/generated"),
        client.V1VolumeMount(name=job_name + "-secret", mount_path="/etc/secret-volume", read_only=True),
        client.V1VolumeMount(name=job_name + "-data", mount_path="/horovod/data")
    ]
    container.command = [
        "/horovod/generated/run.sh"
    ]

    cpu_mode_stub = ""
    if train_mode == "cpu":
        cpu_mode_stub = "ldconfig /usr/local/cuda/lib64/stubs;"

    # TODO: cpu, gpu, 학습 코드
    container.args = [
        "ldconfig /usr/local/cuda/lib64/stubs && mpirun -np {replicas} --hostfile /horovod/generated/hostfile\
            --mca orte_keep_fqdn_hostnames t --allow-run-as-root --display-map --tag-output\
            --timestamp-output sh -c '{cpu_mode_stub} python /horovod/data/train.py'".format(replicas = replicas, cpu_mode_stub = cpu_mode_stub)
    ]

    pod_spec.volumes = [
        client.V1Volume(name=job_name + "-cm",
                        config_map=client.V1ConfigMapVolumeSource(name=job_name,
                                                                  items=[
                                                                      client.V1KeyToPath(key="hostfile.config",
                                                                                         path="hostfile",
                                                                                         mode=438),
                                                                      client.V1KeyToPath(key="master.waitWorkerReady",
                                                                                         path="waitWorkerReady.sh",
                                                                                         mode=365),
                                                                      client.V1KeyToPath(key="master.run",
                                                                                         path="run.sh",
                                                                                         mode=365)

                                                                  ]
                                                                  )),
        client.V1Volume(name=job_name + "-secret",
                        secret=client.V1SecretVolumeSource(secret_name=job_name,
                                                                  default_mode=448,
                                                                  items=[
                                                                      client.V1KeyToPath(key="host-key",
                                                                                         path="id_rsa"),
                                                                      client.V1KeyToPath(key="host-key-pub",
                                                                                         path="authorized_keys")
                                                                  ]

                        )),
        client.V1Volume(name=job_name + "-data",
                        empty_dir=client.V1EmptyDirVolumeSource())
    ]
    pod_spec.containers = [container]

    # init container
    pod_spec.init_containers = [
        client.V1Container(name="wait-workers",
                           image=image,
                           image_pull_policy="IfNotPresent",
                           env=[
                               client.V1EnvVar(name="SSHPORT", value="22"),
                               client.V1EnvVar(name="USESECRETS", value="true")
                           ],
                           command=[
                               "/horovod/generated/waitWorkerReady.sh",
                               # TODO: S3 주소 다시 세팅하기.

                           ],
                           args=[
                               "/horovod/generated/hostfile"
                           ],
                           volume_mounts=[
                                client.V1VolumeMount(name=job_name + "-cm", mount_path="/horovod/generated"),
                                client.V1VolumeMount(name=job_name + "-secret", mount_path="/etc/secret-volume", read_only=True),

                                client.V1VolumeMount(name=job_name + "-data", mount_path="/horovod/data")
                            ]),
        client.V1Container(name="download-data",
                           image=image,
                           image_pull_policy="IfNotPresent",
                           command=[
                               "/bin/bash",
                               "-c"
                           ],
                           args=[
                               "curl http://ywj-horovod.s3.ap-northeast-2.amazonaws.com/horovod/" + username + "/train.py > /horovod/data/train.py"
                           ],
                           volume_mounts=[
                               client.V1VolumeMount(name=job_name + "-data", mount_path="/horovod/data")
                           ])
    ]
    pod_template_spec.spec = pod_spec

    job_spec.template = pod_template_spec
    job.spec = job_spec
    return job


def createService(username, role):
    service_name = username + "-horovod-" + role
    service = client.V1Service()
    service.metadata = client.V1ObjectMeta(name=service_name,
                                           labels={
                                               "app": "horovod",
                                               "user": username,
                                               "role": role
                                           })

    service_spec = client.V1ServiceSpec()
    service_spec.cluster_ip = "None"  # headless Service
    service_spec.selector = {
        "app": "horovod",
        "role": role,
        "user": username
    }
    service_spec.ports = [
        client.V1ServicePort(name="ssh", port=22, target_port=22)
    ]
    service.spec = service_spec
    return service