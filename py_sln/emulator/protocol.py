from typing import Dict

START = 0x68
END = 0x16


def build_ft12_frame(payload: bytes, control: int = 0x53, address: int = 0x01) -> bytes:

    L = 1 + 1 + len(payload)
    frame = bytearray()
    frame.append(START)
    frame.append(L & 0xFF)
    frame.append(L & 0xFF)
    frame.append(START)
    frame.append(control & 0xFF)
    frame.append(address & 0xFF)
    frame.extend(payload)
    cs = sum(frame[4:]) & 0xFF
    frame.append(cs)
    frame.append(END)
    return bytes(frame)


def parse_ft12_frame(buf: bytes) -> Dict:

    if len(buf) < 6:
        raise ValueError("too short")
    if buf[0] != START or buf[3] != START:
        raise ValueError("invalid start bytes")
    L1 = buf[1]
    L2 = buf[2]
    if L1 != L2:
        raise ValueError("length mismatch")
    L = L1
    expected = 4 + L + 2
    if len(buf) != expected:
        raise ValueError(f"frame length {len(buf)} != expected {expected}")
    control = buf[4]
    address = buf[5]
    payload = buf[6 : 6 + (L - 2)]
    cs = buf[6 + (L - 2)]
    calc = sum(buf[4 : 6 + (L - 2)]) & 0xFF
    if cs != calc:
        raise ValueError(f"checksum mismatch {cs} != {calc}")
    if buf[-1] != END:
        raise ValueError("invalid end char")
    return {"control": control, "address": address, "payload": bytes(payload)}
