import boto3
import os
import requests
from urllib.parse import urlparse

def download_from_s3(s3_url: str, dest_path: str):
    """
    Supports s3://bucket/key and https://... style URLs.
    Downloads the image and saves it to dest_path.
    """
    if s3_url.startswith("s3://"):
        parsed = urlparse(s3_url)
        bucket = parsed.netloc
        key = parsed.path.lstrip("/")

        s3 = boto3.client("s3")
        s3.download_file(bucket, key, dest_path)
    elif s3_url.startswith("http"):
        r = requests.get(s3_url)
        r.raise_for_status()
        with open(dest_path, "wb") as f:
            f.write(r.content)
    else:
        raise ValueError(f"Unsupported URL format: {s3_url}")