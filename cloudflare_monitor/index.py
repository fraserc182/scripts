import urllib3
import json
import logging
from datetime import datetime, time, timedelta
import boto3
import os

# initalise and set logging level
logging.basicConfig()
logger = logging.getLogger()
logger.setLevel(logging.INFO)

http = urllib3.PoolManager()

#generic_slack_lambda = os.environ.get("GENERIC_SLACK_LAMBDA")
#slack_channel = os.environ.get("SLACK_ALERT_CHANNEL_NAME")
#queue_url = os.environ.get("SQS_QUEUE_URL")
queue_url = "fraser-test-queue"

sqs = boto3.client('sqs')


def get_unresolved_incidents():
    """
    Get unresolved incidents from Cloudflare
    """
    try:
        r = http.request(
            "GET",
            "https://www.cloudflarestatus.com/api/v2/incidents/unresolved.json",
        )
        if r.status == 200:
            return json.loads(r.data.decode("utf-8"))
        else:
            raise urllib3.exceptions.HTTPError(
                f"HTTP error {r.status}: {r.reason}")

    except (urllib3.exceptions.HTTPError, json.JSONDecodeError) as e:
        logger.exception(
            "Error getting unresolved incidents from Cloudflare: %s", str(e))


def get_status():
    """
    Get current status from Cloudflare
    """
    try:
        r = http.request(
            "GET",
            "https://www.cloudflarestatus.com/api/v2/status.json",
        )
        if r.status == 200:
            return json.loads(r.data.decode("utf-8"))
        else:
            raise urllib3.exceptions.HTTPError(
                f"HTTP error {r.status}: {r.reason}")
    except (urllib3.exceptions.HTTPError, json.JSONDecodeError) as e:
        logger.exception("Error getting status from Cloudflare: %s", str(e))


def notify_slack(incident):


    print(json.dumps({
            "type": "Cloudflare incident alert",
            "message": "Cloudflare is currently experiencing issues, incidents are shown below. For more information visit https://www.cloudflarestatus.com/",
            "channel": slack_channel,
            "additional": {
                "Incident": incident
            }
        }))
    # set lambda & cloudwatch boto3 client
    # lambda_client = boto3.client("lambda")

    # lambda_client.invoke(
    #     FunctionName=generic_slack_lambda,
    #     InvocationType="Event",
    #     Payload=json.dumps({
    #         "type": "Cloudflare incident alert",
    #         "message": "Cloudflare is currently experiencing issues, incidents are shown below. For more information visit https://www.cloudflarestatus.com/",
    #         "channel": slack_channel,
    #         "additional": {
    #             "Incident": incident
    #         }
    #     })
    # )


def check_alert_hours():
    # Check current UTC time
    current_time = datetime.utcnow().time()

    # Define start and end time for alerting (0800-1800 UTC)
    start_time = time(8, 0)  # 0800 UTC
    end_time = time(18, 0)  # 1800 UTC

    # Check if the current time is within the specified time range
    return start_time <= current_time <= end_time


def parse_incidents(incidents):
    # Extract the fields from each dictionary in the list
    extracted_data = []
    for item in incidents:
        extracted_item = {
            'id': item['id'],
            'name': item['name'],
            'status': item['status'],
            'created_at': item['created_at'],
            'updated_at': item['updated_at'],
            'shortlink': item['shortlink']
        }
        extracted_data.append(extracted_item)

    return extracted_data


def process_cloudflare_status():
    #  check if we are within alerting hours
    if not check_alert_hours():
        logger.info("Alert ignored as outside of alerting hours")
        return

    # get the current status of cloudflare and look for major or critical incidents
    status = get_status()
    if status["status"]["indicator"] in ["major", "critical", "minor"]:  # ! for testing
        logger.info("Cloudflare status: %s", status["status"]["description"])
        # retrieve the unresolved incidents
        incidents = get_unresolved_incidents()
        # parse the incidents and split into seperate messages
        msg = parse_incidents(incidents['incidents'])
        msg = [
                {
                    "id": "testid123456",
                    "name": "Cloudflare Pages — Functions live logs issues",
                    "status": "investigating",
                    "created_at": "2023-07-04T14: 35: 20.010Z",
                    "updated_at": "2023-07-04T14: 36: 26.924Z",
                    "shortlink": "https: //stspg.io/ym4z3nth12r2"
                }
            ]

        if msg:
            for inc in msg:
                duplicates = check_for_duplicate_messages(
                    queue_url, 'id', inc['id'])

                if not duplicates:
                    logger.info("no duplicate found, sending to queue")
                    send_message_to_queue(queue_url, inc)
                    logger.info("sending to slack", inc)
                    #notify_slack(inc)
                else:
                    logger.info("duplicate found, not sending to slack")
        else:         
            logger.info("No incidents found")

    else:
        logger.info("Cloudflare is currently operating normally")


def send_message_to_queue(queue_url, message):

    # Define the message attributes
    message_attributes = {
        'id': {'DataType': 'String', 'StringValue': message['id']},
        'name': {'DataType': 'String', 'StringValue': message['name']},
        'status': {'DataType': 'String', 'StringValue': message['status']},
        'created_at': {'DataType': 'String', 'StringValue': message['created_at']},
        'updated_at': {'DataType': 'String', 'StringValue': message['updated_at']},
        'shortlink': {'DataType': 'String', 'StringValue': message['shortlink']}
    }

    # Send the message with attributes to the SQS queue
    response = sqs.send_message(
        QueueUrl=queue_url,
        MessageBody='Your message body',
        MessageAttributes=message_attributes
    )


def check_for_duplicate_messages(queue_url, attribute_name, attribute_value):
    # Set to track processed message IDs
    processed_message_ids = set()

    # List to store duplicate messages
    duplicate_messages = []

    while True:
        # Receive messages from the SQS queue
        response = sqs.receive_message(
            QueueUrl=queue_url,
            AttributeNames=['All'],
            MessageAttributeNames=['All'],
            MaxNumberOfMessages=10,  # Adjust the number of messages to receive
            VisibilityTimeout=60,    # Adjust the visibility timeout as needed
            WaitTimeSeconds=20       # Set the wait time to 0 for non-blocking receive
        )

        # Check if any messages were received
        messages = response.get('Messages', [])
        print(messages)
        if not messages:
            # No messages in the queue, break the loop
            break

        for message in messages:
            # Get the MessageId of the message
            message_id = message['MessageId']

            # Check if the message has already been processed
            if message_id in processed_message_ids:
                duplicate_messages.append(message)

            # Process the message as needed
            attributes = message.get('MessageAttributes', {})
            if attribute_name in attributes and attributes[attribute_name]['StringValue'] == attribute_value:
                duplicate_messages.append(message)

            # Add the message ID to the set of processed messages
            processed_message_ids.add(message_id)

            # # Delete the message from the queue
            # sqs.delete_message(
            #     QueueUrl=queue_url,
            #     ReceiptHandle=message['ReceiptHandle']
            # )

    return duplicate_messages


def lambda_handler(event, context):
    process_cloudflare_status()


lambda_handler(None, None)
