#!/usr/bin/env python3
import os
import datetime
import logging
import json
import time
import boto3

from botocore.exceptions import ClientError, WaiterError

LOG_LEVEL = os.environ.get("LOG_LEVEL", "DEBUG")

logging.basicConfig(format="%(levelname)s:%(message)s")
logger = logging.getLogger()
logger.setLevel(level="INFO")


# set variables
rds_cluster_name = ["customer-replication-db"]
account_id = ["xxx"] # customer account ID

def lambda_handler(event, context):

    load_and_wait_database()
   

def load_and_wait_database():
    dms_client = boto3.client('dms',
                region_name='us-east-2')
    dms_waiter = dms_client.get_waiter('replication_task_stopped')

    task_arns = ['arn:aws:dms:us-east-2:xxx:task:xx',
                    'arn:aws:dms:us-east-2:xxx:task:xxx'
                    ]

    for x in task_arns:
        logger.info((f'Starting replication task {x}'))
        dms_client.start_replication_task(
            ReplicationTaskArn=x,
            StartReplicationTaskType='reload-target'
        )
        time.sleep(90)
        dms_waiter.wait(
            Filters=[
                {
                    'Name': 'replication-task-arn',
                    'Values': [
                        x,
                    ]
                },
            ],
            MaxRecords=100,
            Marker='string',
            WaiterConfig={
                'Delay': 5,
                'MaxAttempts': 30
            }
        )
        logger.info((f'{x} has replicated successfully'))

    return True