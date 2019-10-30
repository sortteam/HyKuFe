import hykufe

print(hykufe.HyKuFeBuilder()
      .setName("test1").setImage("test2")
      .setCpu("test3").setMemory("test4")
      .setGpu("test5").setReplica("test6")
      .build())
