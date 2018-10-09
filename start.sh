PROJECT_PATH=$(pwd)

killProc() {
    server=$1
    ps -fe|grep ${server}|grep -v grep
    if [ $? -ne 0 ]
    then
        echo "${server} not exits"
    else
        echo "killing old ${server}"
        killall ${server}
    fi
}

config() {
    cd ${PROJECT_PATH}
    mkdir -p logs
    chmod 777 logs
}

start() {
    server=$1

    cd ${PROJECT_PATH}/${server}
    echo "start to build ${server}"
    go build .
    echo "starting ${server}"
    nohup ./${server} -c ../conf/$server.yaml > ${PROJECT_PATH}/logs/${server}.log 2>&1 &
}

config
killProc "tcpserver"
start "tcpserver"
killProc "httpserver"
start "httpserver"
