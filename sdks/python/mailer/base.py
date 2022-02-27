from abc import ABCMeta, abstractmethod
from enum import Enum
import json
from typing import List, Optional


class BodyType(Enum):
    PLAIN = "BODY_TYPE_PLAIN"
    HTML = "BODY_TYPE_HTML"


class Client(metaclass=ABCMeta):
    @classmethod
    def __subclasshook__(cls, subclass) -> bool:
        return (
            hasattr(subclass, "_dispatch")
            and callable(subclass.dispatch)
            or NotImplementedError
        )

    @abstractmethod
    def _dispatch(self, path: str, body: str):
        raise NotImplementedError()

    def send(
        self,
        to_email: str,
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        await self._dispatch(
            "/send",
            json.dumps(
                {
                    "to": to_email,
                    "from": from_email,
                    "subject": subject,
                    "body": body,
                    "type": body_type.value,
                    "reply_to": reply_to,
                }
            ),
        )

    def send_batch(
        self,
        to_email: List[str],
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        await self._dispatch(
            "/send/batch",
            json.dumps(
                {
                    "to": to_email,
                    "from": from_email,
                    "subject": subject,
                    "body": body,
                    "type": body_type.value,
                    "reply_to": reply_to,
                }
            ),
        )
