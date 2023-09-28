
use signature::types;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct InputRecord {
    order_id: i32,
    product_id: i32,
    quantity: i32,
    amount: f32,
    shipping: f32,
    tax: f32,
}

#[derive(Serialize, Deserialize, Debug)]
struct OutputRecord {
    order_id: i32,
    product_id: i32,
    quantity: i32,
    total: f32,
}

pub fn scale(ctx: Option<&mut types::Context>) -> Result<Option<types::Context>, Box<dyn std::error::Error>> {
    let unwrapped = ctx.unwrap();

    let unwrapped_input = unwrapped.input.as_ref().unwrap();
    println!("transforming input for topic {} for partition {} and offset {}", unwrapped_input.topic, unwrapped_input.partition, unwrapped_input.offset);

    let input: InputRecord = serde_json::from_slice(unwrapped_input.record.as_slice())?;

    let output = OutputRecord {
        order_id: input.order_id,
        product_id: input.product_id,
        quantity: input.quantity,
        total: input.amount + input.shipping + input.tax,
    };

    let unwrapped_output = unwrapped.output.as_mut().unwrap();
    unwrapped_output.key = unwrapped_input.key.clone();
    unwrapped_output.record = serde_json::to_vec(&output)?;
    unwrapped_output.timestamp = unwrapped_input.timestamp.clone();

    return signature::next(Some(unwrapped));
}
