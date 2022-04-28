use crate::{
    error::Result,
    types::{BodyType, Request},
    Client,
};

pub struct SendBatchBuilder<'s> {
    client: &'s Client,
    to: Vec<&'s str>,
    from: &'s str,
    subject: &'s str,
    body: &'s str,
    body_type: Option<BodyType>,
    reply_to: Option<&'s str>,
}

impl<'s> SendBatchBuilder<'s> {
    pub(crate) fn new(
        client: &'s Client,
        from: &'s str,
        subject: &'s str,
        body: &'s str,
    ) -> SendBatchBuilder<'s> {
        SendBatchBuilder {
            client,
            to: Vec::new(),
            from,
            subject,
            body,
            body_type: None,
            reply_to: None,
        }
    }

    /// Set all the recipients of the message
    pub fn recipients(mut self, recipients: Vec<&'s str>) -> SendBatchBuilder<'s> {
        self.to = recipients;
        self
    }

    /// Add a new recipient to the message
    pub fn to(mut self, to: &'s str) -> SendBatchBuilder<'s> {
        self.to.push(to);
        self
    }

    /// Set the type of the message body
    pub fn body_type(mut self, body_type: BodyType) -> SendBatchBuilder<'s> {
        self.body_type = Some(body_type);
        self
    }

    /// Set the reply to of the message
    pub fn reply_to(mut self, reply_to: &'s str) -> SendBatchBuilder<'s> {
        self.reply_to = Some(reply_to);
        self
    }

    /// Send the request
    pub async fn send(self) -> Result<()> {
        let req = Request::new(
            self.to,
            self.from,
            self.subject,
            self.body,
            self.body_type,
            self.reply_to,
        );

        self.client.dispatch("/send/batch", req).await
    }
}
