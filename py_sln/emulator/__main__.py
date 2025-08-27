from .cli import parse_args
from .logging_config import init_logging
from .server import Emulator
import asyncio
import logging


def main():
    args = parse_args()
    init_logging()
    em = Emulator(args.host, args.port, args.proto, args.mode)
    try:
        asyncio.run(em.start())
    except KeyboardInterrupt:
        logging.info("Stopped by user")


if __name__ == "__main__":
    main()
