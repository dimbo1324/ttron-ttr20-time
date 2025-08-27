import datetime


def make_text_response() -> bytes:
    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    return f"TIME:{now}\n".encode("utf-8")


def make_ft12_time_payload() -> bytes:
    return f"TIME:{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}".encode(
        "ascii"
    )
