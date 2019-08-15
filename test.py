from sort import TrainingClusterManager

train_script = "tensorflow2_mnist.py"

tm = TrainingClusterManager("ywj", token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6Inl3ai1zYS10b2tlbi1sbnh3bSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJ5d2otc2EiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI1YzRjMGYzZi1iOWYwLTExZTktYTNhYy1mYTE2M2UwOTllM2YiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDp5d2otc2EifQ.dxalsvH1kTLbD6FULnp3K6ukxebqpkxAM3myq2dYJIETKL1VGfGOPxVmY5jEjiXE_Lgb7uKHu7kckXDoVieeOljGxUPE1wGRNlicWHm2BKu57QQ8wQimaERJTrxHbxuh6d-U90lU28Yg4Y-7BmGunp2VJeBAZ2ajsNRMKw-u1l38glaOQlKFtNby94KzcaNA0jPsTXOFKKf7UMmJiUXhlJ1Pf-RhqHZ72jV4ZXr1OT9hpFJBEqUP9C4iC5k8prc36IyUr6-9mDwWlFV6VhPvtGt9DkoAF7DVYvj0MrZYrreVtFfZTuv5LMdgce6hjWI3e0wLkMnBL8o0K57Z1PC3UQ",
                     replicas=1,
                     kube_master_ip="203.254.143.253:8080",
                     entry_point=train_script,
                     image="horovod/horovod:0.16.4-tf1.14.0-torch1.1.0-mxnet1.4.1-py3.6",
                     is_host_network=True,
                     ssh_port="10111"
                     )

#tm.sendDataToS3()
tm.deleteAllObject()

tm.runTrain()

# model, weight = tm.runTrain()
# model.save()
# weight.save()