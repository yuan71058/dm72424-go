# -*- coding: utf-8 -*-
"""把 onnx 模型转换为低 opset 版本 (大漠插件需要 opset 12-14)"""
import onnx
from onnx import version_converter
from pathlib import Path

src = Path("yolov5s_single.onnx")
dst = Path("yolov5s_dm.onnx")

print(f"加载: {src}")
model = onnx.load(str(src), load_external_data=True)
print(f"原始 opset: {[(o.domain, o.version) for o in model.opset_import]}")

target_opset = 12
print(f"转换为 opset {target_opset} ...")
try:
    converted = version_converter.convert_version(model, target_opset)
    print(f"转换后 opset: {[(o.domain, o.version) for o in converted.opset_import]}")
    onnx.save_model(converted, str(dst), save_as_external_data=False)
    print(f"✓ 保存: {dst} ({dst.stat().st_size} bytes)")
except Exception as e:
    print(f"✗ version_converter 失败: {e}")
    print("尝试手动设置 opset ...")
    # 手动降级: 修改 opset 版本号 (不保证所有算子兼容, 但大漠可能能加载)
    for opset in model.opset_import:
        opset.version = target_opset
    model.ir_version = 7  # 兼容较老 runtime
    onnx.save_model(model, str(dst), save_as_external_data=False)
    print(f"✓ 手动降级保存: {dst} ({dst.stat().st_size} bytes)")
