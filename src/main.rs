use axum_server::tls_rustls::RustlsConfig;
use mci::{app, config::Config};
use std::{net::SocketAddr, path::PathBuf};
use tracing::{info, warn};
use tracing_subscriber::EnvFilter;

#[tokio::main]
async fn main() {
    let config = Config::from_env().expect("Failed to load configuration from environment");

    tracing_subscriber::fmt()
        .with_target(false)
        .with_env_filter(EnvFilter::new(&config.log_level))
        .init();

    let app = app();
    let addr: SocketAddr = config
        .address
        .parse()
        .expect("Invalid address format in MCI_ADDRESS");
    let handle = axum_server::Handle::new();
    let shutdown_handle = handle.clone();

    tokio::spawn(async move {
        tokio::signal::ctrl_c()
            .await
            .expect("Failed to listen for Ctrl+C");

        info!("Shutdown signal received. Closing server gracefully...");

        shutdown_handle.graceful_shutdown(Some(std::time::Duration::from_secs(30)));
    });

    if let (Some(cert_path), Some(key_path)) = (config.cert_path, config.key_path) {
        let tls_config =
            RustlsConfig::from_pem_file(PathBuf::from(&cert_path), PathBuf::from(&key_path))
                .await
                .expect("Failed to load TLS certificates from provided paths");

        info!("Server listening on {}", addr);

        axum_server::bind_rustls(addr, tls_config)
            .handle(handle)
            .serve(app.into_make_service())
            .await
            .unwrap();
    } else {
        warn!("TLS certificates not provided. Starting insecure HTTP server.");
        info!("Server listening on {}", addr);

        axum_server::bind(addr)
            .handle(handle)
            .serve(app.into_make_service())
            .await
            .unwrap();
    }
}
