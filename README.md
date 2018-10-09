Simple User-Management system
===========================

Requirements
------------
* bash
* Go(v1.10.3+)
* Mysql(v5.7+)
* Redis(v4.0+)

Installation
------------
```$xslt
    #makesure database and redis is correctlly installed and started
    sudo sysctl -w kern.ipc.somaxconn=2048
    sudo sysctl -w kern.maxfiles=12288
    ulimit -n 10000

    #clone code
    cd $GOPATH
    git clone https://git.garena.com/jinhua.ouyang/entry_task.git src

    #setup environment
    cd src
    sh install.sh
```

Start
------------
```$xslt
    cd src
    sh start.sh

    #logs
    refer to src/logs
```

Performance testing
------------
install apache tools
```
    sudo apt-get install apache2-utils
```

test with web client
```
    http://192.168.33.10:8080/static/login.html
    username/pwd : test/test
```

test with curl
```
    curl -d "username=test&passwd=098f6bcd4621d373cade4e832627b4f6" "http://localhost:8080/login"
```

performance test with ab
add "username=test&passwd=098f6bcd4621d373cade4e832627b4f6" into login.txt, remember use *set noeol; w ++bin* to strip endline if using vim

for login test
```
    ab -n 50000 -c 2000 -T 'application/x-www-form-urlencoded' -p login.txt "http://localhost:8080/login"
```
for random login test
```
    ab -n 50000 -c 2000 -T 'application/x-www-form-urlencoded' -p empty.txt "http://localhost:8080/randlogin"
```
