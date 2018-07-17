# Image 学习



可以尝试
```bash
docker image ls
```
查看你有的镜像，结果可能如下所示

```
REPOSITORY                   TAG                 IMAGE ID            CREATED             SIZE
web                          latest              a650b6f161f8        4 days ago          1.06GB
python                       3.6                 d330010a503a        11 days ago         912MB
gosqlserverdemo_web          latest              56f0929033d3        4 weeks ago         420MB
gosqlserverdemo              p1                  2447287b7ce5        4 weeks ago         420MB
gosqlserverdemo              t1                  7b7b29967786        4 weeks ago         419MB
golang                       1.10-alpine         c4b5d89b27f4        5 weeks ago         376MB
```

你可以选择ls的参数删除镜像等。下面我们来构建一下镜像



## 分层的构建

在这里，我们会运行一个 Nginx 容器：

```bash
docker run --name webserver -d -p 80:80 nginx
```

你在你的电脑上可以访问 localhost:80.

```bash
docker exec -it webserver bash
root@31ed98fd5559:/# echo '<h1>GC Course: Web & Docker</h1>' > /usr/share/nginx/html/index.html
root@31ed98fd5559:/# exit
exit
```

我们可以运行 [Docker Diff](https://docs.docker.com/engine/reference/commandline/diff/#parent-command)

```
docker diff webserver
C /root
A /root/.bash_history
C /run
A /run/nginx.pid
C /usr/share/nginx/html/index.html
C /var/cache/nginx
A /var/cache/nginx/client_temp
A /var/cache/nginx/fastcgi_temp
A /var/cache/nginx/proxy_temp
A /var/cache/nginx/scgi_temp
A /var/cache/nginx/uwsgi_temp
```

这里A 表示添加， C表示修改，我们可以commit

```bash
docker commit --author "mwish <1506118561@qq.com>" --message "修改了默认网页" webserver nginx:v2

sha256:c60b73064e691fc1e834669fb5e33baa717d772e229c126b5cbb01f6fd259164
```

可以看到对应的条目

```
nginx                        v2                  c60b73064e69        7 seconds ago       108MB
```

你也可以运行这个容器，达到你想要的目的

参考这个 https://zhuanlan.zhihu.com/p/31744232 你可以看到你萌的操作

History 可以看到你的修改。在每一层中，你的修改会被标出来，Docker 分层构建你的镜像

```
IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
c60b73064e69        5 minutes ago       nginx -g daemon off;                            115B                修改了默认网页
e4e6d42c70b3        12 months ago       /bin/sh -c #(nop)  CMD ["nginx" "-g" "daemon…   0B                  
```

## Dockerfile 构建镜像

对于以上的构建，我们可以看到，以上的操作的构建是不透明、难重复而且细节很诡异的，我们可以用 Dockerfile 构建镜像。以上的 Dockerfile 编写如下：

```dockerfile
FROM nginx
RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html
```

在目录下运行 `docker build` 指令，可以构建出我们想要的镜像

```bash
docker build -t nginx:v3 .
```



```
Step 1/2 : FROM nginx


 ---> e4e6d42c70b3

Step 2/2 : RUN echo '<h1>Hello, Docker!</h1>' > /usr/share/nginx/html/index.html


 ---> Running in 3ffeb870777b

Removing intermediate container 3ffeb870777b

 ---> 929e734948a3

Successfully built 929e734948a3

Successfully tagged nginx:v3

'nginx:v3 Dockerfile: Images/nginx-v3/Dockerfile' has been deployed successfully.
```

这里是我们得到的输出，可以看到，`FROM`指定基础镜像 , `RUN`表示执行命令。我们可以看到 FROM构建了基本的惊险，RUN执行运行之后，产生了新的景象，随之构建成功。

##  基本的应用讲解

这里我搭建了一个小应用，详见代码。

```dockerfile
FROM golang:1.10-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache git

ENV SRC_DIR=/go/src/github.com/mapleFU/Images/dockerfile-learn

WORKDIR $SRC_DIR

COPY ./ $SRC_DIR

EXPOSE 8080 8080

RUN go build $SRC_DIR/main.go

CMD ["./main"]
```

`RUN` 这里是运行指令。 && 是因为 RUN 每一句都会构建

```
docker build -t learn-docker:v1 .

Sending build context to Docker daemon  29.18MB
Step 1/8 : FROM golang:1.10-alpine
 ---> c4b5d89b27f4
Step 2/8 : RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories &&     apk add --no-cache git
 ---> Using cache
 ---> 75edd7ec5adc
Step 3/8 : ENV SRC_DIR=/go/src/github.com/mapleFU/Images/dockerfile-learn
 ---> Using cache
 ---> 039793765ecd
Step 4/8 : WORKDIR $SRC_DIR
 ---> Using cache
 ---> cb96a3bbb005
Step 5/8 : COPY ./ $SRC_DIR
 ---> c466e11f7e10
Step 6/8 : EXPOSE 8080 8080
 ---> Running in 25f3d044ec06
Removing intermediate container 25f3d044ec06
 ---> 56990965e6eb
Step 7/8 : RUN go build $SRC_DIR/main.go
 ---> Running in b51e75786be6
Removing intermediate container b51e75786be6
 ---> 4ad8a2c99db4
Step 8/8 : CMD ["./main"]
 ---> Running in 64da48e6efad
Removing intermediate container 64da48e6efad
 ---> a33b3ab1dc0a
Successfully built a33b3ab1dc0a
Successfully tagged learn-docker:v1
```

我们可以根据这个看到构建镜像的过程，我们可以运行它

```bash
docker run learn-docker:v1
```

我们可以用*httpie*进行测试

```
http -v --json POST localhost:8000/login username=mwish password=mwish
```

结果我们会发现，测试失败：

```
http: error: ConnectionError: HTTPConnectionPool(host='localhost', port=8080): Max retries exceeded with url: /login (Caused by NewConnectionError('<urllib3.connection.HTTPConnection object at 0x109447f28>: Failed to establish a new connection: [Errno 61] Connection refused',)) while doing POST request to URL: http://localhost:8080/login
```

我们需要在run的时候映射端口（虽然这是下一节的内容），再操作：

```
docker run -p 8080:8080 learn-docker:v1
```



```http
POST /login HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 42
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/0.9.9

{
    "password": "mwish",
    "username": "mwish"
}

```

```http
HTTP/1.1 200 OK
Content-Length: 207
Content-Type: application/json; charset=utf-8
Date: Sun, 15 Jul 2018 13:59:08 GMT

{
    "code": 200,
    "expire": "2018-07-15T14:59:08Z",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzE2NjY3NDgsImlkIjoibXdpc2giLCJvcmlnX2lhdCI6MTUzMTY2MzE0OH0.d6iGRlCjARsY4BPEvNTrtLRWLOIbi_yj-p5ufpqQAiw"
}
```

```
http -f GET localhost:8008/auth/hello "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzE2NjY3NDgsImlkIjoibXdpc2giLCJvcmlnX2lhdCI6MTUzMTY2MzE0OH0.d6iGRlCjARsY4BPEvNTrtLRWLOIbi_yj-p5ufpqQAiw" "Content-Type: application/json"
```

```
http -f GET localhost:8080/auth/hello "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzE2NjY3NDgsImlkIjoibXdpc2giLCJvcmlnX2lhdCI6MTUzMTY2MzE0OH0.d6iGRlCjARsY4BPEvNTrtLRWLOIbi_yj-p5ufpqQAiw" "Content-Type: application/json"
```

```http
HTTP/1.1 200 OK
Content-Length: 40
Content-Type: application/json; charset=utf-8
Date: Sun, 15 Jul 2018 14:27:57 GMT

{
    "text": "Hello World.",
    "userID": "mwish"
}

```

