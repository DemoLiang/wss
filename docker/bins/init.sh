#!/bin/bash
service ssh start
/root/landfw &
supervisord -n
