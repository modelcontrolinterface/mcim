use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
};
use thiserror::Error;

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Error, Debug)]
pub enum Error {
    #[error("an unexpected error occurred")]
    Anyhow(#[from] anyhow::Error),
}

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let (status, error_message) = match self {
            Error::Anyhow(_) => (
                StatusCode::INTERNAL_SERVER_ERROR,
                "an unexpected error occurred".to_string(),
            ),
        };

        (status, error_message).into_response()
    }
}
