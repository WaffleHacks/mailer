from enum import Enum


class BodyType(Enum):
    PLAIN = "BODY_TYPE_PLAIN"
    HTML = "BODY_TYPE_HTML"


class InvalidArgumentException(Exception):
    """
    A passed argument is invalid
    """
