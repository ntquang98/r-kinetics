import os
import tempfile
from typing import Dict

import numpy as np
import tensorflow as tf
from collections import defaultdict
from utils import download_from_s3

MODEL_PATH = "./models/student_distilled_V_0.3.tflite"  # Update path if needed

BINARY_LABELS = ['non-snake', 'snake']
VENOM_LABELS  = ['Unknown', 'non-venomous', 'venomous']

class InferenceEngine:
    def __init__(self, model_path: str = MODEL_PATH):
        if not os.path.isfile(model_path):
            raise FileNotFoundError(f"Model file not found: {model_path}")

        self.interp = tf.lite.Interpreter(model_path=model_path)
        self.interp.allocate_tensors()
        self.inp = self.interp.get_input_details()[0]

        dim2idx = defaultdict(list)
        for o in self.interp.get_output_details():
            dim2idx[o['shape'][-1]].append(o['index'])

        if len(dim2idx[2]) != 1 or len(dim2idx[3]) != 1:
            raise RuntimeError(f"Expected one 2-class and one 3-class output, got {dict(dim2idx)}")

        self.head_to_index = {
            'snake_out': dim2idx[2][0],
            'venom_out': dim2idx[3][0],
        }

    def predict(self, image_path: str) -> Dict:
        img = tf.io.read_file(image_path)
        img = tf.image.decode_jpeg(img, channels=3)
        img = tf.image.resize(img, (480,480)) / 255.0
        arr = img.numpy().astype(np.float32)[None, ...]

        dt = self.inp['dtype']
        if dt in (np.int8, np.uint8):
            scale, zp = self.inp['quantization']
            arr = np.round(arr / scale + zp).astype(dt)

        self.interp.set_tensor(self.inp['index'], arr)
        self.interp.invoke()

        results = {}
        for head, idx in self.head_to_index.items():
            od = next(o for o in self.interp.get_output_details() if o['index'] == idx)
            raw = self.interp.get_tensor(idx)[0]
            if od['dtype'] in (np.int8, np.uint8):
                scale, zp = od['quantization']
                probs = (raw.astype(np.float32) - zp) * scale
            else:
                probs = raw.astype(np.float32)

            cid = int(np.argmax(probs))
            conf = float(probs[cid])
            label = BINARY_LABELS[cid] if head == 'snake_out' else VENOM_LABELS[cid]
            results[head] = {
                "label": label,
                "confidence": conf
            }

        return results

    def predict_from_s3_url(self, s3_url: str) -> Dict:
        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=True) as tmp:
            download_from_s3(s3_url, tmp.name)
            return self.predict(tmp.name)
