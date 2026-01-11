use axum::Router;
use tower_http::trace::TraceLayer;

pub mod error;
pub mod model;
pub mod services;
pub mod web;

pub fn app() -> Router {
    Router::new()
        .merge(web::routes::routes())
        .layer(TraceLayer::new_for_http())
}
