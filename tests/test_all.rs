use anyhow::Result;
use axum::{
    body::{to_bytes, Body},
    http::{Request, StatusCode},
};
use mci::{app, model::Message};
use tower::ServiceExt;

#[tokio::test]
async fn test_json_message_route() -> Result<()> {
    let response = app()
        .oneshot(Request::builder().uri("/json").body(Body::empty()).unwrap())
        .await
        .unwrap();

    assert_eq!(response.status(), StatusCode::OK);

    let body = to_bytes(response.into_body(), usize::MAX).await?;
    let message: Message = serde_json::from_slice(&body)?;

    assert_eq!(message.message, "Hello from the API");

    Ok(())
}
