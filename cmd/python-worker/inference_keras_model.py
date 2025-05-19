import os
import tempfile
from typing import Dict
import numpy as np
from tensorflow.keras.models import load_model
from utils import download_from_s3

MODEL_PATH = "./models/version_3a_with_2heads.h5"
THRESHOLD = 0.75

BINARY_LABELS = ['non-snake', 'snake']
VENOM_LABELS  = ['Unknown', 'non-venomous', 'venomous']

class InferenceEngine:
    def __init__(self, model_path: str = MODEL_PATH):
        if not os.path.isfile(model_path):
            raise FileNotFoundError(f"Model file not found: {model_path}")
        self.model = load_model(model_path)

    def predict(self, image_path: str) -> Dict:
        import tensorflow as tf
        img = tf.io.read_file(image_path)
        img = tf.image.decode_jpeg(img, channels=3)
        img = tf.image.resize(img, (480, 480)) / 255.0
        arr = img.numpy()[None, ...].astype(np.float32)

        outputs = self.model.predict(arr)

        # Expect 2 heads: [snake_pred, venom_pred]
        snake_probs = outputs[0][0]
        venom_probs = outputs[1][0]

        snake_idx = int(np.argmax(snake_probs))
        snake_conf = float(snake_probs[snake_idx])
        snake_label = BINARY_LABELS[snake_idx]

        venom_idx = int(np.argmax(venom_probs))
        venom_conf = float(venom_probs[venom_idx])
        venom_label = VENOM_LABELS[venom_idx]

        result = {
            "snake_out": {
                "label": snake_label,
                "confidence": snake_conf
            },
            "venom_out": {
                "label": venom_label,
                "confidence": venom_conf
            }
        }

        if snake_conf < THRESHOLD:
            result["snake_out"]["label"] = "uncertain"
            result["snake_out"]["note"] = f"Confidence below threshold ({THRESHOLD})"

        print("=======RUN", result)

        return result

    def predict_from_s3_url(self, s3_url: str) -> Dict:
        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=True) as tmp:
            download_from_s3(s3_url, tmp.name)
            return self.predict(tmp.name)
