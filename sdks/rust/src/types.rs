use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize)]
pub enum Format {
    #[serde(rename = "BODY_TYPE_PLAIN")]
    Plain,
    #[serde(rename = "BODY_TYPE_HTML")]
    Html,
}

#[derive(Debug, Serialize)]
pub(crate) struct Request<'r, T> {
    pub to: T,
    pub from: &'r str,
    pub subject: &'r str,
    pub body: &'r str,
    pub format: Format,
    pub reply_to: Option<&'r str>,
}

impl<'r, T> Request<'r, T> {
    pub(crate) fn new(
        to: T,
        from: &'r str,
        subject: &'r str,
        body: &'r str,
        format: Option<Format>,
        reply_to: Option<&'r str>,
    ) -> Self {
        Request {
            to,
            from,
            subject,
            body,
            format: format.unwrap_or(Format::Plain),
            reply_to,
        }
    }
}

#[derive(Debug, Deserialize)]
pub(crate) struct Response {
    pub message: String,
}
