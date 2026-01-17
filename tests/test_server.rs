use http_body_util::Empty;
use hyper::{body::Bytes, Request, Version};
use hyper_util::rt::tokio::TokioIo;
use mci::{app, config::Config, db, s3, AppState};
use std::net::SocketAddr;
use tokio::net::TcpStream;

mod utils;

async fn spawn_app() -> SocketAddr {
    std::thread::sleep(std::time::Duration::from_secs(5));

    utils::generate_certs().await;

    std::env::set_var("MCI_ADDRESS", "127.0.0.1:0");
    std::env::set_var(
        "MCI_DATABASE_URL",
        "postgres://postgres:postgres@localhost:5432/mci",
    );
    std::env::set_var("MCI_S3_URL", "http://localhost:8333");
    std::env::set_var("MCI_S3_ACCESS_KEY", "none");
    std::env::set_var("MCI_S3_SECRET_KEY", "none");

    let config = Config::from_env().unwrap();
    let db_pool = db::create_pool(&config.database_url);
    let s3_client =
        s3::create_s3_client(&config.s3_url, &config.s3_access_key, &config.s3_secret_key).await;

    db::init_db(&db_pool).await.unwrap();

    let app = app(AppState { db_pool, s3_client });
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let addr = listener.local_addr().unwrap();

    tokio::spawn(async move {
        axum::serve(listener, app.into_make_service())
            .await
            .unwrap();
    });

    addr
}

#[tokio::test]
async fn test_http1() {
    let addr = spawn_app().await;
    let stream = TokioIo::new(TcpStream::connect(addr).await.unwrap());
    let (mut sender, conn) = hyper::client::conn::http1::handshake(stream).await.unwrap();

    tokio::spawn(conn);

    let request = Request::builder()
        .version(Version::HTTP_11)
        .uri(format!("http://{}/test", addr))
        .body(Empty::<Bytes>::new())
        .unwrap();
    let response = sender.send_request(request).await.unwrap();

    assert!(response.status().is_success());
}

#[tokio::test]
async fn test_http2() {
    let addr = spawn_app().await;
    let stream = TcpStream::connect(addr).await.unwrap();
    let (mut client, h2) = h2::client::handshake(stream).await.unwrap();

    tokio::spawn(async move {
        if let Err(e) = h2.await {
            println!("GOT ERR: {:?}", e);
        }
    });

    let request = Request::builder()
        .version(Version::HTTP_2)
        .uri(format!("http://{}/test", addr))
        .body(())
        .unwrap();
    let (response_future, _) = client.send_request(request, true).unwrap();
    let response = response_future.await.unwrap();

    assert!(response.status().is_success());
}
