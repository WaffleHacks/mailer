from requests import Session
from typing import Dict, List, Optional
from yarl import URL

from .shared import BodyType, InvalidArgumentException


class Client(object):
    """
    A blocking mailer client for sending messages
    """

    def __init__(self, server: str):
        self.base_url = URL(server)

        self.session = Session()
        self.session.headers.update({"Content-Type": "application/json"})

    def close(self):
        """
        Close the connection to the mailer
        """
        self.session.close()

    def _dispatch(self, path: str, body: Dict[str, str]):
        url = self.base_url.with_path(path)
        response = self.session.post(str(url), json=body)

        if response.status_code == 400:
            body = response.json()
            raise InvalidArgumentException(body["message"])

    def send(
        self,
        to_email: str,
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        """
        Send a single email
        :param to_email: the address of the recipient
        :param from_email: the address of the sender in RFC 5322
        :param subject: the email subject
        :param body: the message body
        :param body_type: the content type of the body
        :param reply_to: an optional email to reply to
        """

        self._dispatch(
            "/send",
            {
                "to": to_email,
                "from": from_email,
                "subject": subject,
                "body": body,
                "type": body_type.value,
                "reply_to": reply_to,
            },
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
        """
        Send an email to many recipients
        :param to_email: the addresses of the recipients
        :param from_email: the address of the sender in RFC 5322
        :param subject: the email subject
        :param body: the message body
        :param body_type: the content type of the body
        :param reply_to: an optional email to reply to
        """

        self._dispatch(
            "/send/batch",
            {
                "to": to_email,
                "from": from_email,
                "subject": subject,
                "body": body,
                "type": body_type.value,
                "reply_to": reply_to,
            },
        )

    def send_template(
        self,
        to: Dict[str, Dict[str, str]],
        from_email: str,
        subject: str,
        body: str,
        body_type: BodyType = BodyType.PLAIN,
        reply_to: Optional[str] = None,
    ):
        """
        Send a templated email to many recipients
        :param to: the addresses of the recipients in RFC 5322 format with their associated contexts
        :param from_email: the address of the sender in RFC 5322 format
        :param subject: the email subject
        :param body: the message body template
        :param body_type: the content type of the body
        :param reply_to: an optional email to reply to
        """
        # Transform the to contexts
        prepared_to = {}
        for key, context in to.items():
            prepared_to[key] = {
                "key": list(context.keys()),
                "value": list(context.values()),
            }

        self._dispatch(
            "/send/template",
            {
                "to": prepared_to,
                "from": from_email,
                "subject": subject,
                "body": body,
                "type": body_type.value,
                "reply_to": reply_to,
            },
        )
