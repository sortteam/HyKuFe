from sort import TrainingClusterManager

train_script = "tensorflow2_mnist.py"

tm = TrainingClusterManager("ywj",
                     replicas=2,
                     entry_point=train_script,
                     docker_image="sort/sort:0.16.4-tf1.14.0-torch1.1.0-mxnet1.4.1-py3.6",
                     is_host_network=False
                     )

#tm.sendDataToS3()
tm.deleteAllObject()

tm.runTrain()

# model, weight = tm.runTrain()
# model.save()
# weight.save()