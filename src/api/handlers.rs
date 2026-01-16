use axum::Json;
use crate::{
    errors::test::Result,
    domains::test::Message,
};

pub async fn json_message() -> Result<Json<Message>> {
    let message = Message {
        message: "Hello from the API".to_string(),
    };

    Ok(Json(message))
}
