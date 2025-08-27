import logging
import sys


def init_logging(logfile: str = "emulator.log", level: int = logging.INFO):
    root = logging.getLogger()
    if root.handlers:
        for h in list(root.handlers):
            root.removeHandler(h)
    handlers = [
        logging.FileHandler(logfile, encoding="utf-8"),
        logging.StreamHandler(sys.stdout),
    ]
    logging.basicConfig(
        level=level, format="%(asctime)s %(message)s", handlers=handlers
    )
