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
account_id = ["xxxxxxx"]

def _make_snapshot_shareable(snapshotObject, rds_client, attributeName='restore', valuesToAdd=[]):
    waiter = rds_client.get_waiter("db_cluster_snapshot_available")
    try:
        snapShotIdentifier = snapshotObject['DBClusterSnapshot']['DBClusterSnapshotIdentifier']
    except KeyError:
        logger.error("Issue getting SnapshotIdentifier from snapshot object. Please ensure it's the correct format: {}".format(
            snapshotObject))
        return

    logger.info("waiting for snapshot to become available")
    try:
        waiter.wait(
            DBClusterSnapshotIdentifier=snapShotIdentifier,
            SnapshotType="manual"
        )
    except WaiterError as e:
        logger.error("Waiter Error: {}".format(e))

    try:
        res = rds_client.modify_db_cluster_snapshot_attribute(
            DBClusterSnapshotIdentifier=snapShotIdentifier,
            AttributeName=attributeName,
            ValuesToAdd=valuesToAdd
        )
    except ClientError as e:
        logger.error("error modifying snapshot. ValuesToAdd={}, snapshotIdentifier={}, error: {} ".format(
            valuesToAdd, snapShotIdentifier, e))
    else:
        logger.info(
            "shareable snapshot successfully generated for {}".format(valuesToAdd))
        logger.info("{}".format(json.dumps(res)))


def _generate_rds_snapshot(clusterIdentifier, rds_client):
    """This is ultimately for customer snapshot sharing"""
    timestamp = datetime.datetime.now()
    timestamp = timestamp.strftime("%Y-%m-%d-%H%M%S")
    snapshot_identifier = f"{clusterIdentifier}-{timestamp}"

    snapshot_tags = [
        {
            "Key": "instance",
            "Value": clusterIdentifier
        }
    ]

    logger.info("Preparing snapshot with the following details: {}".format(
        json.dumps(
            {
                "snapshot_identifier": snapshot_identifier
            }
        )
    ))

    try:
        response = rds_client.create_db_cluster_snapshot(
            DBClusterSnapshotIdentifier=snapshot_identifier,
            DBClusterIdentifier=clusterIdentifier,
            Tags=snapshot_tags
        )
    except ClientError as e:
        logger.error("Error generating RDS Snapshot for DB Instance {}. Error is: {}".format(
            clusterIdentifier, e))
        return
    else:
        return response


def lambda_handler(event, context):

    if load_and_wait_database():
        rds_client = boto3.client("rds",
                region_name='us-east-2')
                
        for x in rds_cluster_name:
            res = _generate_rds_snapshot(
                clusterIdentifier=x, rds_client=rds_client)
            _make_snapshot_shareable(res, rds_client, valuesToAdd=account_id)

    else:
        logger.error("Issue replicating database")
   
def load_and_wait_database():
    dms_client = boto3.client('dms',
                region_name='us-east-2')
    dms_waiter = dms_client.get_waiter('replication_task_stopped')

    task_arns = ['arn:aws:dms:us-east-2:xxxx:task:12344567',
                    'arn:aws:dms:us-east-2:xxxx:task:1234567']

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
                'Delay': 60,
                'MaxAttempts': 120
            }
        )
        logger.info((f'{x} has replicated successfully'))

    return True


event = ""
context = ""
lambda_handler(event,context)