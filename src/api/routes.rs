use crate::api::handlers;
use axum::{routing::get, Router};

pub fn routes() -> Router {
    Router::new().route("/test", get(handlers::json_message))
}
