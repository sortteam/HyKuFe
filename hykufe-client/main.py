import hykufe

# hykufe.HyKuFeBuilder()\
#     .setName("test1").setImage("test2")\
#     .setCPU("test3").setMemory("test4")\
#     .setGPU("test5").setReplica("test6")\
#     .build().writeYamlFile("test.yaml")

hykufe.HyKuFeBuilder()..build('access_key', 'secret_key').uploadFileToS3('main.py')

