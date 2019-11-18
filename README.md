# HyKuFe
## 개요

학습 코드와 데이터만 설정하면 자동으로 Horovod 분산 학습 환경을 설정하고 학습을 해주는 파이썬 코드입니다.



Defined Resource
1. Job - Master 파드 1개 실행
2. StatefulSet - Worker 파드 N개 실행
3. Service - Master, Worker 파드에 대한 각각 Headless Service 생성
4. ConfigMap - 파드 초기 설정(SSH 설정 등), Master가 Worker 상태를 확인하는 쉘스크립트 포함
5. Secret - SSH Private, Public Key 등 정보



## 사용법
1. 먼저 TrainingClusterManager 객체를 생성합니다. 생성자의 인자 목록은 다음과 같습니다.

Argument | 설명 | Required | Default Value
---|---|---|---
kube_master_ip| 쿠버네티스 마스터 API 주소 | True | 
token| 쿠버네티스 마스터 API에 접근하기 위한 토큰 | True |
username | 사용자 이름 | True |
entry_point | 학습 스크립트 | True |
train_data | 학습 데이터 | False |
test_data | 테스트 데이터 | False |
train_node_num | 분산 학습을 진행할 노드 개수 | False | `2`
namespace | | False | `default`
image | Horovod 학습 파드에 띄울 이미지 | False | `uber/horovod:0.12.1-tf1.8.0-py3.5`
is_host_network | 호스트 네트워크 사용 유무. 네트워크로 인한 오버헤드 차이가 큼.[True, False] | False | `False`
ssh_port | 호스트 네트워크로 설정했을 시 사용할 ssh port | False | `22`

2. runTrain함수를 호출합니다.

## BUG
1. 아직 Node Affinity가 적용되지 않아 파드들이 의도한 대로(노드 당 1개) 띄지 않음. 호스트 네트워크 사용시 문제 발생.

## TODO
- 학습 종료 후 학습 결과와 파라미터 값들을 저장하는 기능
  - 어디에 저장? 사용자 Jupyter 환경에? 혹은 대시보드를 만든다면 Sagemaker처럼
- Node Affinity, Resource Request 옵션 추가(추후 Hybrid 환경에서 학습을 진행한다면 Cluster Affinity도 고려해야 할듯. kubefed 참고.
- SSH Key이 현재는 고정으로 파일에 있는데 학습 실행 시 키 파일을 생성하도록 수정
- Persistent Volume Claim 추가하여 학습 코드, 학습 데이터 다운받아 저장하도록 수정(일단.) (현재는 EmptyDir로 각 파드 생성시 각자 다운받도록 되어있음... ^^)
  - Master Container에서 initContainer에서 처리하게 하면 될듯.
- S3 Public Access 해제하고 Access Key나 뭐 기타 방법으로 보안 강화
- 사용자 학습 코드 부분에서 데이터를 어떻게 받고 학습 결과물을 어디에 저장할 껀지 설정하는 argument가 필요할듯.

 

## 고민할 부분
- 현재는 kubenetes master API에 직접 접근하고 python client를 이용하여 클러스터를 제어하고 있음.
  - Final -> Custom Resource Controller로 컨트롤 하는게 Best 일듯
- 사용자가 학습을 취소할 때 (즉 Stop 버튼을 누를 때) 코드가 실행이 중단되므로 만들어진 학습 Object 제어권이 소멸되버림.
  - 해결하려면 사용자 레이어가 아닌 쿠버네티스 클러스터 내에서 Resource들을 
- 학습 도중 cluster의 node에서 fail이 일어난 경우
- 사용자가 여러개의 학습을 돌리고 싶은 경우 Schedular에 대한 고민도 필요할 듯.
