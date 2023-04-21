from enum import Enum


class Format(Enum):
    PLAIN = "PLAIN"
    HTML = "HTML"


class InvalidArgumentException(Exception):
    """
    A passed argument is invalid
    """
