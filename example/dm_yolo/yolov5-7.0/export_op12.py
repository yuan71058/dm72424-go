# -*- coding: utf-8 -*-
"""
使用 torch.onnx.export 手动导出 opset 12 的 yolov5s onnx
大漠插件需要 opset 12-14, PyTorch 2.x 默认导出 opset 18 不兼容.
本脚本直接调用底层 torch.onnx.export 绕过 yolov5 export.py 的默认行为.
"""
import sys
import torch
from pathlib import Path

# 加载 yolov5 模型
sys.path.insert(0, str(Path(".").resolve()))
from models.common import DetectMultiBackend
from utils.general import check_img_size

weights = "yolov5s.pt"
device = "cpu"
model = DetectMultiBackend(weights, device=device)
stride = model.stride
imgsz = check_img_size(640, s=stride)
model.eval()

# 创建 dummy input
dummy = torch.zeros(1, 3, imgsz, imgsz)

# 更新模型 forward 以适配 onnx 导出
from models.yolo import Detect
m = model.model.model[-1] if hasattr(model.model, 'model') else model.model[-1]
if isinstance(m, Detect):
    m.inplace = False
    m.export = True

out_path = "yolov5s_op12.onnx"
print(f"导出 opset 12 onnx -> {out_path}")

with torch.no_grad():
    torch.onnx.export(
        model.model if hasattr(model, 'model') else model,
        dummy,
        out_path,
        opset_version=12,
        input_names=["images"],
        output_names=["output0"],
        do_constant_folding=True,
        dynamic_axes=None,  # 静态形状, 大漠可能更兼容
        dynamo=False,  # 强制使用旧版 TorchScript trace 导出器, 尊重 opset_version
    )

print(f"✓ 导出完成: {out_path}")

# simplify (大漠官方推荐 --simplify)
import onnx
try:
    from onnxsim import simplify
    m2 = onnx.load(out_path)
    m2_simp, check = simplify(m2)
    if check:
        onnx.save(m2_simp, out_path)
        print(f"✓ simplify 完成 (opset 保持 {m2_simp.opset_import[0].version})")
    else:
        print("⚠ simplify 校验失败, 保留原始模型")
except Exception as e:
    print(f"⚠ simplify 跳过: {e}")

# 添加 metadata (关键! 大漠 7.2336+ 从 onnx metadata 读取类名, 不再依赖 .class 文件)
# 与 yolov5 官方 export.py 的 metadata 格式保持一致
m2 = onnx.load(out_path)
stride_val = int(model.stride.max()) if hasattr(model.stride, 'max') else int(model.stride)
d = {'stride': stride_val, 'names': model.names}
for k, v in d.items():
    meta = m2.metadata_props.add()
    meta.key, meta.value = k, str(v)
onnx.save(m2, out_path)
print(f"✓ metadata 已写入: stride={d['stride']}, names={len(d['names'])} 个类名")

# 验证
m2 = onnx.load(out_path, load_external_data=False)
print(f"  IR: {m2.ir_version}, Opset: {[(o.domain, o.version) for o in m2.opset_import]}")
print(f"  大小: {Path(out_path).stat().st_size} bytes")
print(f"  输入: {[(i.name, [d.dim_value for d in i.type.tensor_type.shape.dim]) for i in m2.graph.input]}")
print(f"  输出: {[(o.name, [d.dim_value for d in o.type.tensor_type.shape.dim]) for o in m2.graph.output]}")
print(f"  Metadata: {[(p.key, p.value[:80]) for p in m2.metadata_props]}")

# onnxruntime 验证
try:
    import onnxruntime as ort
    sess = ort.InferenceSession(out_path)
    print(f"  onnxruntime 验证通过")
except Exception as e:
    print(f"  onnxruntime 验证失败: {e}")
