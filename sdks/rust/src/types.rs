use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize)]
pub enum BodyType {
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
    pub body_type: BodyType,
    pub reply_to: Option<&'r str>,
}

impl<'r, T> Request<'r, T> {
    pub(crate) fn new(
        to: T,
        from: &'r str,
        subject: &'r str,
        body: &'r str,
        body_type: Option<BodyType>,
        reply_to: Option<&'r str>,
    ) -> Self {
        Request {
            to,
            from,
            subject,
            body,
            body_type: body_type.unwrap_or(BodyType::Plain),
            reply_to,
        }
    }
}

#[derive(Debug, Serialize)]
pub(crate) struct TemplateContext<'c> {
    key: Vec<&'c str>,
    value: Vec<&'c str>,
}

impl<'c> From<HashMap<&'c str, &'c str>> for TemplateContext<'c> {
    fn from(h: HashMap<&'c str, &'c str>) -> Self {
        TemplateContext {
            key: h.keys().copied().collect(),
            value: h.values().copied().collect(),
        }
    }
}

#[derive(Debug, Deserialize)]
pub(crate) struct Response {
    pub message: String,
}
