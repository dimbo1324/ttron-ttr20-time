import argparse


def parse_args():
    p = argparse.ArgumentParser(
        description="Simple Teleport K-104 emulator (FT1.2-like)"
    )
    p.add_argument("--host", default="0.0.0.0", help="bind host")
    p.add_argument("--port", type=int, default=9000, help="port to bind")
    p.add_argument(
        "--proto", choices=["tcp", "udp"], default="tcp", help="transport protocol"
    )
    p.add_argument(
        "--mode",
        choices=["text", "binary"],
        default="text",
        help="behaviour mode (unused)",
    )
    return p.parse_args()
