from aiohttp import ClientSession

from .base import Client


class AsyncClient(Client):
    """
    An async mailer client for sending messages
    """

    def __init__(self, server: str):
        base_url = server
        if not base_url.endswith("/"):
            base_url += "/"

        self.session = ClientSession(base_url)

    def close(self):
        """
        Close the connection to the mailer
        """
        self.session.close()

    async def _dispatch(self, path: str, body: str):
        response = await self.session.post(
            path, data=body, headers={"Content-Type": "application/json"}
        )
        print(await response.json())
