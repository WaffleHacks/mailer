use reqwest::Error as ReqwestError;
use thiserror::Error;

pub(crate) type Result<T> = std::result::Result<T, Error>;

#[derive(Debug, Error)]
pub enum Error {
    #[error("failed to connect to server")]
    Connection,
    #[error("invalid argument: {0}")]
    InvalidArgument(String),
    #[error("unable to parse response")]
    Parse,
    #[error("request timed out")]
    Timeout,
    #[error("an unknown error occurred: {0}")]
    Unknown(String),
}

impl From<ReqwestError> for Error {
    fn from(e: ReqwestError) -> Self {
        if e.is_body() || e.is_decode() {
            Error::Parse
        } else if e.is_timeout() {
            Error::Timeout
        } else if e.is_connect() {
            Error::Connection
        } else {
            Error::Unknown(format!("{e}"))
        }
    }
}
