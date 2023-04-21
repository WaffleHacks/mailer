use crate::{
    error::Result,
    types::{Format, Request},
    Client,
};
use std::collections::HashMap;

#[derive(Debug)]
pub struct SendTemplateBuilder<'s> {
    client: &'s Client,
    to: HashMap<&'s str, HashMap<&'s str, &'s str>>,
    from: &'s str,
    subject: &'s str,
    body: &'s str,
    format: Option<Format>,
    reply_to: Option<&'s str>,
}

impl<'s> SendTemplateBuilder<'s> {
    pub(crate) fn new(
        client: &'s Client,
        from: &'s str,
        subject: &'s str,
        body: &'s str,
    ) -> SendTemplateBuilder<'s> {
        SendTemplateBuilder {
            client,
            to: HashMap::new(),
            from,
            subject,
            body,
            format: None,
            reply_to: None,
        }
    }

    /// Set the type of the message body
    pub fn format(mut self, format: Format) -> SendTemplateBuilder<'s> {
        self.format = Some(format);
        self
    }

    /// Set the reply to of the message
    pub fn reply_to(mut self, reply_to: &'s str) -> SendTemplateBuilder<'s> {
        self.reply_to = Some(reply_to);
        self
    }

    /// Add a new recipient to the message
    pub fn to(
        mut self,
        to: &'s str,
        context: Option<HashMap<&'s str, &'s str>>,
    ) -> SendTemplateBuilder<'s> {
        self.to.insert(to, context.unwrap_or_default());
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

        self.client.dispatch("/send/template", req).await
    }
}
