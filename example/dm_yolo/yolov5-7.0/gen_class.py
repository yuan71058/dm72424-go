# -*- coding: utf-8 -*-
"""从 coco.yaml 生成 yolov5s.class 文件 (大漠 YOLO 要求与 onnx 同名同目录)"""
import yaml
from pathlib import Path

yaml_path = Path("data/coco.yaml")
out_path = Path("yolov5s.class")

with open(yaml_path, "r", encoding="utf-8") as f:
    d = yaml.safe_load(f)

names = [d["names"][i] for i in sorted(d["names"])]

with open(out_path, "w", encoding="utf-8") as f:
    f.write("\n".join(names) + "\n")

print(f"生成 {out_path}, 共 {len(names)} 个类")
print("前5个:", names[:5])
