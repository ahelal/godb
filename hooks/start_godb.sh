#!/bin/sh

JOB_NAME=godb
JOBS_DIR=/var/masnae/jobs
JOB_DIR="${JOBS_DIR}/${JOB_NAME}"
LOG_DIR="/var/masnae/log/${JOB_NAME}"

exec ${JOB_DIR}/pkg/bin/goDB-linux --config ${JOB_DIR}/config/config.yml >>  $LOG_DIR/${JOB_NAME}.stdout.log 2>> $LOG_DIR/${JOB_NAME}.stderr.log
