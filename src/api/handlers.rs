use crate::{domains::test::Message, errors::test::Result};
use axum::Json;

pub async fn json_message() -> Result<Json<Message>> {
    let message = Message {
        message: "Hello from the API".to_string(),
    };

    Ok(Json(message))
}
