import logging
from ..protocol import START, parse_ft12_frame, build_ft12_frame
from ..responses import make_text_response, make_ft12_time_payload
import asyncio


async def handle_tcp(
    reader: asyncio.StreamReader, writer: asyncio.StreamWriter, server=None
):
    addr = writer.get_extra_info("peername")
    logging.info("TCP conn from %s", addr)
    try:
        buf = bytearray()
        while True:
            chunk = await reader.read(1024)
            if not chunk:
                break
            buf.extend(chunk)
            logging.info("Received (%s): %r", addr, chunk)
            while True:
                if len(buf) >= 1 and buf[0] == START:
                    if len(buf) < 6:
                        break
                    L1 = buf[1]
                    L2 = buf[2]
                    if L1 != L2:
                        logging.warning(
                            "Length bytes mismatch (%d != %d), resync", L1, L2
                        )
                        buf.pop(0)
                        continue
                    L = L1
                    expected = 4 + L + 2
                    if len(buf) < expected:
                        break
                    frame = bytes(buf[:expected])
                    try:
                        info = parse_ft12_frame(frame)
                        logging.info(
                            "Parsed FT1.2 request from %s: control=0x%02X addr=0x%02X payload=%r",
                            addr,
                            info["control"],
                            info["address"],
                            info["payload"],
                        )
                        resp_payload = make_ft12_time_payload()
                        resp = build_ft12_frame(
                            resp_payload, control=0x73, address=info["address"]
                        )
                    except Exception as e:
                        logging.warning("FT1.2 parse failed: %s; fallback to text", e)
                        resp = make_text_response()
                    writer.write(resp)
                    await writer.drain()
                    logging.info("Sent (%s): %r", addr, resp)
                    buf = buf[expected:]
                    continue
                else:
                    if b"\n" in buf:
                        idx = buf.find(b"\n")
                        line = bytes(buf[: idx + 1])
                        try:
                            s = line.decode("utf-8", errors="ignore").strip()
                        except Exception:
                            s = ""
                        logging.info("Text line from %s: %r", addr, s)
                        if s.upper() == "GETTIME":
                            resp = make_text_response()
                        else:
                            resp = make_text_response()
                        writer.write(resp)
                        await writer.drain()
                        logging.info("Sent (%s): %r", addr, resp)
                        buf = buf[idx + 1 :]
                        continue
                    else:
                        break
    except Exception:
        logging.exception("TCP handler error")
    finally:
        logging.info("Connection closed %s", addr)
        try:
            writer.close()
            await writer.wait_closed()
        except Exception:
            pass
