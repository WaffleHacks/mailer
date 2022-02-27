use reqwest::{
    header::{HeaderMap, HeaderValue},
    Client as HttpClient, StatusCode, Url,
};
use std::time::Duration;

mod error;
mod types;

pub use error::Error;
use error::Result;
pub use types::BodyType;
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

    async fn handle_response(resp: reqwest::Response) -> Result<()> {
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
    pub async fn send<T, F, S, B>(
        &self,
        to: T,
        from: F,
        subject: S,
        body: B,
        body_type: Option<BodyType>,
        reply_to: Option<&str>,
    ) -> Result<()>
    where
        T: AsRef<str>,
        F: AsRef<str>,
        S: AsRef<str>,
        B: AsRef<str>,
    {
        let resp = self
            .client
            .post(self.base_url.join("/send").unwrap())
            .json(&Request::new(
                to.as_ref(),
                from.as_ref(),
                subject.as_ref(),
                body.as_ref(),
                body_type,
                reply_to,
            ))
            .send()
            .await?;

        Self::handle_response(resp).await
    }

    /// Send an email to many recipients
    pub async fn send_batch<'s, T, F, S, B>(
        &self,
        to: T,
        from: F,
        subject: S,
        body: B,
        body_type: Option<BodyType>,
        reply_to: Option<&str>,
    ) -> Result<()>
    where
        T: AsRef<[&'s str]>,
        F: AsRef<str>,
        S: AsRef<str>,
        B: AsRef<str>,
    {
        let resp = self
            .client
            .post(self.base_url.join("/send/batch").unwrap())
            .json(&Request::new(
                to.as_ref(),
                from.as_ref(),
                subject.as_ref(),
                body.as_ref(),
                body_type,
                reply_to,
            ))
            .send()
            .await?;

        Self::handle_response(resp).await
    }
}
