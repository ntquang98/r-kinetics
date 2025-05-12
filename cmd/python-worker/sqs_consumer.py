import json
import os
import time
import traceback

import boto3
import requests
from inference import InferenceEngine

print("Using AWS creds:", os.environ.get("AWS_ACCESS_KEY_ID"), os.environ.get("AWS_REGION"))

QUEUE_URL = os.environ.get("SQS_QUEUE_URL")
POLL_INTERVAL = 5

sqs = boto3.client("sqs", region_name=os.environ["AWS_REGION"])
engine = InferenceEngine()

def process_message(msg_body):
    try:
        print("========msg body", msg_body)
        data = json.loads(msg_body)
        job_id = data["job_id"]
        s3_url = data["file_url"]
        callback_url = data["callback_url"]

        if job_id == "" or s3_url == "" or callback_url == "":
            print("job is not valid", job_id, s3_url, callback_url)
            return

        result = engine.predict_from_s3_url(s3_url)

        print("result", result)

        resultStr = json.dumps(result)

        resp = requests.put(callback_url, json={
            "id": job_id,
            "result": resultStr
        })

        print("Response status:", resp.status_code)
        print("Response body:", resp.text)  

    except Exception as e:
        print("[ERROR]", str(e))
        traceback.print_exc()
        try:
            requests.put(data.get("callback_url", ""), json={
                "id": data.get("job_id", "unknown"),
                "message": "failed" + str(e)
            })
        except:
            pass

def poll():
    while True:
        resp = sqs.receive_message(
            QueueUrl=QUEUE_URL,
            MaxNumberOfMessages=5,
            WaitTimeSeconds=10
        )
        for msg in resp.get("Messages", []):
            receipt = msg["ReceiptHandle"]
            process_message(msg["Body"])
            sqs.delete_message(QueueUrl=QUEUE_URL, ReceiptHandle=receipt)
        time.sleep(POLL_INTERVAL)

if __name__ == '__main__':
    print("[SQS Worker] Polling for jobs...")
    poll()
