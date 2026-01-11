use axum_server::tls_rustls::RustlsConfig;
use mci::app;
use std::{net::SocketAddr, path::PathBuf};
use tracing::info;
use tracing_subscriber::EnvFilter;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_target(false)
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let app = app();
    let addr = SocketAddr::from(([127, 0, 0, 1], 8080));

    let config = RustlsConfig::from_pem_file(
        PathBuf::from("certs/cert.pem"),
        PathBuf::from("certs/key.pem"),
    )
    .await
    .unwrap();

    info!("->> LISTENING on {addr} with TLS and HTTP/2\n");
    axum_server::bind_rustls(addr, config)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
