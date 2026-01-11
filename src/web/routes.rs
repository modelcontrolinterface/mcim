use crate::web::handlers;
use axum::{routing::get, Router};

pub fn routes() -> Router {
    Router::new().route("/json", get(handlers::json_message))
}
