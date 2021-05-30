# nicping可以同时对多个ip进行ping，并生成易读的结果

## 使用方法：
```
Usage of ./nicping:
  -c int
        执行ping的并发数量。默认是10
  -hosts string
        存放ip或域名的文件，文件中每行存放一个ip或域名。默认文件是当前目录的hosts.txt。如果没有文件，可以直接指定一个ip
  -p int
        每次ping发送的包个数。默认是2个
```