import asyncio
import logging
from .handlers import udp as udp_handler_module, tcp as tcp_handler_module


class Emulator:
    def __init__(self, host: str, port: int, proto: str, mode: str):
        self.host = host
        self.port = port
        self.proto = proto.lower()
        self.mode = mode.lower()

    async def start(self):
        logging.info(
            "Starting emulator on %s:%d proto=%s mode=%s",
            self.host,
            self.port,
            self.proto,
            self.mode,
        )
        if self.proto == "tcp":
            server = await asyncio.start_server(
                lambda r, w: tcp_handler_module.handle_tcp(r, w, self),
                host=self.host,
                port=self.port,
            )
            addrs = ", ".join(str(sock.getsockname()) for sock in server.sockets or [])
            logging.info("TCP server listening on %s", addrs)
            async with server:
                await server.serve_forever()
        else:
            loop = asyncio.get_running_loop()
            transport, protocol = await loop.create_datagram_endpoint(
                lambda: udp_handler_module.UDPHandler(self),
                local_addr=(self.host, self.port),
            )
            logging.info("UDP server listening on %s:%d", self.host, self.port)
            try:
                while True:
                    await asyncio.sleep(3600)
            finally:
                transport.close()
