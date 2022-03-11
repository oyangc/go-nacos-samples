#!/bin/bash

### BEGIN INIT INFO
# Provides:             paas-nacos
# Required-Start:       $syslog $remote_fs
# Required-Stop:        $syslog $remote_fs
# Should-Start:         $local_fs
# Should-Stop:          $local_fs
# Default-Start:        2 3 4 5
# Default-Stop:         0 1 6
# Short-Description:    paas-nacos
# Description:          paas-nacos
### END INIT INFO


# update
NAME=paas-nacos
DIR_HOME=/opt/paas/paas-nacos
EXEC=/opt/jdk/bin/java
PIDFILE=${DIR_HOME}/run.pid
chmod +x ${DIR_HOME}/go-nacos

echo "${DIR_HOME}"


log_out() {
  msg="$*"
  echo $msg
  echo "`date '+%Y-%m-%d %H:%M:%S'` ${msg}" >>${DIR_HOME}/run.out 2>&1
}

start_server(){
   msg="Starting $NAME ..."
   ${DIR_HOME}/go-nacos
   log_out ${msg}
}

stop_server(){
    msg="Stop $NAME ..."
    log_out ${msg}
}

case "$1" in
    start)
        start_server
        ;;
    stop)
        stop_server
        ;;
    restart)
        stop_server
        start_server
        ;;
    *)
        stop_server
        start_server
        echo "Usage: $NAME {start|stop|restart}" >&2
        ;;
esac
