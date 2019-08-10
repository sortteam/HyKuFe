# PythonTrainingManager
## 개요

Resource
1. Job - Master 파드 1개 실행
2. StatefulSet - Worker 파드 N개 실행
3. Service - Master, Worker 파드에 대한 각각 Headless Service 생성
4. ConfigMap - 파드 초기 설정(SSH 설정 등), Master가 Worker 상태를 확인하는 쉘스크립트 포함.
5. Secret - SSH Private, Public Key

현재는 학습 코드만 S3에 전송하고 사용받는 식으로 되어있음. ( 학습 데이터 구현 X)

기본 이미지는 uber/horovod:0.12.1-tf1.8.0-py3.5로 되어있으나 
horovod/horovod:0.16.4-tf1.14.0-torch1.1.0-mxnet1.4.1-py3.6 이것도 되는거 확인

 

## TODO
- 학습 종료 후 학습 결과와 파라미터 값들을 저장하는 기능
  - 어디에 저장? 사용자 Jupyter 환경에? 혹은 대시보드를 만든다면 Sagemaker처럼
- FL을 위한 Tensor.js로 모델 컨버팅 기능. (내가 해야하나 사용자가 해야하나?)
- FL을 위한 모델 CDN 업로드.(일단은 S3가 편할듯.)
- Node Affinity, Resource Request 옵션 추가(추후 Hybrid 환경에서 학습을 진행한다면 Cluster Affinity도 고려해야 할듯. kubefed 참고.
- SSH Key이 현재는 고정으로 파일에 있는데 학습 실행 시 키 파일을 생성하도록 수정
- Persistent Volume Claim 추가하여 학습 코드, 학습 데이터 다운받아 저장하도록 수정(일단.) (현재는 EmptyDir로 각 파드 생성시 각자 다운받도록 되어있음... ^^)
  - Master Container에서 initContainer에서 처리하게 하면 될듯.
- S3 Public Access 해제하고 Access Key나 뭐 기타 방법으로 보안 강화
- Pod 상태 모니터링 및 Log 출력 버그 고치기
- 사용자 학습 코드 부분에서 데이터를 어떻게 받고 학습 결과물을 어디에 저장할 껀지 설정하는 argument가 필요할듯.

 

## 고민할 부분
- 현재는 kubenetes master API에 직접 접근하고 python client를 이용하여 클러스터를 제어하지만 딱히 좋은 방법같진 않음
  - Final. Custom Resource Controller를 만드는게 Best 일듯(But, 난이도 ⭐ ⭐ ⭐ ⭐ ⭐ )
- 사용자가 학습을 취소할 때 (즉 Stop 버튼을 누를 때) 코드가 실행이 중단되므로 만들어진 학습 Object 제어권이 소멸되버림.
  - 해결하려면 사용자 레이어가 아닌 큐브 클러스터에서 Resource들을 
- 학습 도중 cluster의 node에서 fail이 일어난 경우
- 사용자가 여러개의 학습을 돌리고 싶은 경우 Scedular에 대한 고민도 필요할 듯.
