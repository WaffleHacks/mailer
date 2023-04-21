use crate::{
    error::Result,
    types::{Format, Request},
    Client,
};

pub struct SendBuilder<'s> {
    client: &'s Client,
    to: &'s str,
    from: &'s str,
    subject: &'s str,
    body: &'s str,
    format: Option<Format>,
    reply_to: Option<&'s str>,
}

impl<'s> SendBuilder<'s> {
    pub(crate) fn new(
        client: &'s Client,
        to: &'s str,
        from: &'s str,
        subject: &'s str,
        body: &'s str,
    ) -> SendBuilder<'s> {
        SendBuilder {
            client,
            to,
            from,
            subject,
            body,
            format: None,
            reply_to: None,
        }
    }

    /// Set the type of the message body
    pub fn format(mut self, format: Format) -> SendBuilder<'s> {
        self.format = Some(format);
        self
    }

    /// Set the reply to of the message
    pub fn reply_to(mut self, reply_to: &'s str) -> SendBuilder<'s> {
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
            self.format,
            self.reply_to,
        );

        self.client.dispatch("/send", req).await
    }
}
