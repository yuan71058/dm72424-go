# -*- coding: utf-8 -*-
"""把带外部数据的 onnx 重新保存为单文件 onnx (大漠插件要求单文件)"""
import onnx
from pathlib import Path

src = Path("yolov5s.onnx")
dst = Path("yolov5s_single.onnx")

print(f"加载: {src}")
model = onnx.load(str(src), load_external_data=True)

print(f"保存为单文件: {dst}")
onnx.save_model(model, str(dst), save_as_external_data=False)

print(f"✓ 完成, 大小: {dst.stat().st_size} bytes")
