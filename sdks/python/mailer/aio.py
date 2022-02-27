from aiohttp import ClientSession
import json
from typing import List, Optional

from .base import Client, BodyType


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

    async def send(
        self,
        to_email: str,
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        data = json.dumps(
            {
                "to": to_email,
                "from": from_email,
                "subject": subject,
                "body": body,
                "type": body_type.value,
                "replyTo": reply_to,
            }
        )
        response = await self.session.post(
            "/send", data=data, headers={"Content-Type": "application/json"}
        )
        print(await response.json())

    async def send_batch(
        self,
        to_email: List[str],
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        data = json.dumps(
            {
                "to": to_email,
                "from": from_email,
                "subject": subject,
                "body": body,
                "type": body_type.value,
                "replyTo": reply_to,
            }
        )
        response = await self.session.post(
            "/send", data=data, headers={"Content-Type": "application/json"}
        )
        print(await response.json())
