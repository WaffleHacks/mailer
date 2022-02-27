from abc import ABCMeta, abstractmethod
from enum import Enum
from typing import List, Optional


class BodyType(Enum):
    PLAIN = 1
    HTML = 2


class Client(metaclass=ABCMeta):
    @classmethod
    def __subclasshook__(cls, subclass) -> bool:
        has_send = hasattr(subclass, "send") and callable(subclass.send)
        has_send_batch = hasattr(subclass, "send_batch") and callable(
            subclass.send_batch
        )
        return has_send and has_send_batch

    @abstractmethod
    def send(
        self,
        to_email: str,
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        raise NotImplementedError()

    @abstractmethod
    def send_batch(
        self,
        to_email: List[str],
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        raise NotImplementedError()
