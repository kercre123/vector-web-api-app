#!/bin/bash

ssh root@$1 "systemctl stop vic-custom-web"
scp ./build/custom-web root@$1:/bin/
rsync -avr var/www/* root@$1:/var/www/
ssh root@$1 "systemctl start vic-custom-web"
ssh root@$1
