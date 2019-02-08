#!/bin/sh

#
# JOB_NAME
# RUN_DIR=/var/vcap/sys/run/web_ui
# LOG_DIR=/var/vcap/sys/log/web_ui

JOB_NAME=godb
JOBS_DIR=/var/masnae/jobs
JOB_DIR="${JOBS_DIR}/${JOB_NAME}"
LOG_DIR="/var/masnae/log/${JOB_NAME}"
# PIDFILE=${RUN_DIR}/pid

# mkdir -p $RUN_DIR $LOG_DIR
# chown -R vcap:vcap $RUN_DIR $LOG_DIR

# echo $$ > $PIDFILE

# cd /var/vcap/packages/ardo_app

# export PATH=/var/vcap/packages/ruby_1.9.3/bin:$PATH

# exec /var/vcap/packages/ruby_1.9.3/bin/bundle exec \
#       rackup -p <%= p('port') %> \
#       >>  $LOG_DIR/web_ui.stdout.log \
#       2>> $LOG_DIR/web_ui.stderr.log

exec ${JOB_DIR}/pkg/bin/godb --config ${JOB_DIR}/config/config.yml >>  $LOG_DIR/${JOB_NAME}.stdout.log 2>> $LOG_DIR/${JOB_NAME}.stderr.log
