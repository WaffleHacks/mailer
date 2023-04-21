use reqwest::{
    header::{HeaderMap, HeaderValue},
    Client as HttpClient, StatusCode, Url,
};
use serde::Serialize;
use std::time::Duration;

mod builders;
mod error;
mod types;

use builders::*;
pub use error::Error;
use error::Result;
pub use types::Format;
use types::{Request, Response};

/// A client for the WaffleHacks mailer
#[derive(Debug)]
pub struct Client {
    client: HttpClient,
    base_url: Url,
}

impl Client {
    /// Create a new mailer client
    pub fn new(base_url: Url) -> Self {
        // Create the default headers
        let mut headers = HeaderMap::new();
        headers.insert("Content-Type", HeaderValue::from_static("application/json"));

        // This shouldn't ever return an error
        let client = HttpClient::builder()
            .default_headers(headers)
            .timeout(Duration::from_secs(10))
            .build()
            .unwrap();

        Client { client, base_url }
    }

    /// Send a request to the server
    pub(crate) async fn dispatch<T>(&self, path: &str, req: Request<'_, T>) -> Result<()>
    where
        T: Serialize,
    {
        let resp = self
            .client
            .post(self.base_url.join(path).unwrap())
            .json(&req)
            .send()
            .await?;

        let status = resp.status();
        if status == StatusCode::OK {
            Ok(())
        } else {
            let body: Response = resp.json().await?;
            if status == StatusCode::BAD_REQUEST {
                Err(Error::InvalidArgument(body.message))
            } else {
                Err(Error::Unknown(body.message))
            }
        }
    }

    /// Send a single email
    pub async fn send<'s>(
        &'s self,
        to: &'s str,
        from: &'s str,
        subject: &'s str,
        body: &'s str,
    ) -> SendBuilder<'s> {
        SendBuilder::new(self, to, from, subject, body)
    }

    /// Send an email to many recipients
    pub async fn send_batch<'s>(
        &'s self,
        from: &'s str,
        subject: &'s str,
        body: &'s str,
    ) -> SendBatchBuilder<'s> {
        SendBatchBuilder::new(self, from, subject, body)
    }

    /// Send a templated email to many recipients
    pub async fn send_template<'s>(
        &'s self,
        from: &'s str,
        subject: &'s str,
        body: &'s str,
    ) -> SendTemplateBuilder<'s> {
        SendTemplateBuilder::new(self, from, subject, body)
    }
}
