import struct, sys

p = r'E:\SRC\go-dm72424\example\dm_yolo\dm_yolo.exe'
d = open(p, 'rb').read()
if d[:2] != b'MZ':
    print('NOT A PE FILE'); sys.exit(0)
e_lfanew = struct.unpack('<I', d[0x3c:0x40])[0]
sig = d[e_lfanew:e_lfanew+4]
machine = struct.unpack('<H', d[e_lfanew+4:e_lfanew+6])[0]
# PE32 / PE32+ optional header magic at e_lfanew+4+24 (2 bytes): 0x10b=PE32, 0x20b=PE32+
magic = struct.unpack('<H', d[e_lfanew+4+24:e_lfanew+4+26])[0]
print('file size =', len(d))
print('e_lfanew  =', hex(e_lfanew))
print('sig       =', sig)
print('machine   =', hex(machine), '(0x14c=386, 0x8664=amd64)')
print('opt magic =', hex(magic), '(0x10b=PE32/32bit, 0x20b=PE32+/64bit)')
