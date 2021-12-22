#!/bin/bash

## Script for starting, stopping, restarting, and checking the current status of fisync
## This script can be run as follows: ./restart.sh <start/stop/status>

export BEEGO_RUNMODE=prod

PWD=`pwd`
exec="fisync"
chmod +x ${exec}
version=`./fisync -v`
binpath=${PWD}

echo "
                __  __       _     _ 
 ___  ___ __ _ / _|/ _| ___ | | __| |
/ __|/ __/ _` | |_| |_ / _ \| |/ _` |
\__ \ (_| (_| |  _|  _| (_) | | (_| |
|___/\___\__,_|_| |_|  \___/|_|\__,_|

${version}     
============================================================================================        
"
pid=`ps -ef | grep "${binpath}/$exec" | grep -v "grep" | awk '{print $2}'`
option=$1

case $option in
        start)
                if [[ $pid -gt 1 ]]; then
                        echo "${exec}"
                        echo "************WARN: process has exist************"
                        echo ""
                        exit 1
                fi
                echo "==========starting...=========="
                cd ${binpath}
                mkdir -p ${logDir}
                nohup ${binpath}/${exec}  >> ${logDir}/${exec}.log 2>&1 &
                cd -
                echo "==========started=============="
                echo ""
        ;;
        stop)
                if [[ $pid -gt 1 ]]; then
                        echo "==============stop...================"
                        kill -15 $pid
                        echo "[$exec] [$pid]"
                        echo "==============stoped================="
                else
                        echo "************WARN: pid not found !!!************"
                fi
                echo ""
        ;;
        restart)
                if [[ $pid -gt 1 ]]; then
                        echo "============restarting============"
                        chmod +x ${binpath}/${exec}
                        kill -HUP $pid
                        echo "[$exec] [$pid]"
                        echo "============restarted============="
                else
                        echo "[$exec]"
                        echo "************WARN: not started************"
                fi
                echo ""
        ;;
        *)
        echo "$0:usage: [start|stop|restart]"
        exit 1
        ;;
esac