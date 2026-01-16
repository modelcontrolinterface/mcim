use axum::Router;
use tower_http::trace::TraceLayer;

pub mod api;
pub mod config;
pub mod domains;
pub mod errors;

pub fn app() -> Router {
    Router::new()
        .merge(api::routes::routes())
        .layer(TraceLayer::new_for_http())
}
