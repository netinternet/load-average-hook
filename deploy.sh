#!/bin/bash
host=$1
if [[ "$host" =~ ^(([1-9]?[0-9]|1[0-9][0-9]|2([0-4][0-9]|5[0-5]))\.){3}([1-9]?[0-9]|1[0-9][0-9]|2([0-4][0-9]|5[0-5]))$ ]]; then
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build .
    mv load-average-hook loadhook
    ssh -t root@$host "systemctl stop loadhook"
    scp ./loadhook root@$host:/usr/bin/loadhook 
    scp ./loadhook.conf root@$host:/tmp/loadhook.conf 
    ssh -t root@$host "[ ! -f /etc/loadhook.conf ] && cp -r /tmp/loadhook.conf /etc/loadhook.conf || rm -rf /tmp/loadhook.conf"
    scp ./loadhook.service root@$host:/tmp/loadhook.service
    ssh -t root@$host "[ ! -f /etc/systemd/system/loadhook.service ] && cp -r /tmp/loadhook.service /etc/systemd/system/loadhook.service || rm -rf /tmp/loadhook.service"
    ssh -t root@$host "systemctl daemon-reload && systemctl restart loadhook && systemctl enable loadhook"
else
  echo "Example: bash deploy.sh 127.0.0.1"
  exit 1 
fi

