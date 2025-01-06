use serde_json::{from_str, Value};
use tungstenite::{connect, Message};

const QUERY: &str = r#"{"jsonrpc": "2.0", "method": "subscribe", "id": 0, "params": {"query": "tm.event = 'Tx' AND message.action = '/canine_chain.storage.MsgBuyStorage'"}}"#;

fn main() {
    let mut cmd_str = String::new();
    let mut cmd: i16;
    loop {
        println!("1: Subscribe 3: Exit"); // no unsubscribe, not dealing with concurrency
        std::io::stdin().read_line(&mut cmd_str).expect("no input");
        cmd = cmd_str.trim().parse().expect("not integer");
        if cmd == 1 || cmd == 3 {
            break;
        }
    }

    if cmd == 1 {
        let (mut socket, _) =
            connect("ws://localhost:26657/websocket").expect("can't connect");
        socket.send(Message::Text(QUERY.into())).unwrap();

        let mut msg: String;
        let mut msg_json: Value;
        loop {
            msg = socket.read().expect("disconnect").to_string();
            msg_json = from_str(&msg).expect("invalid json");
            if msg.contains("MsgBuyStorage") {
                let now: chrono::DateTime<chrono::Utc> = chrono::Utc::now();
                let bytes = msg_json["result"]["events"]
                    ["buy_storage.bytes_bought"][0]
                    .as_str()
                    .unwrap()
                    .parse::<f64>()
                    .unwrap();
                let hours = msg_json["result"]["events"]
                    ["buy_storage.hours_bought"][0]
                    .as_str()
                    .unwrap()
                    .parse::<f64>()
                    .unwrap();
                println!(
                    "{} | {:.2} gb {:.2} days {} buyer {} receiver {} tx",
                    now.format("%F %T").to_string(),
                    bytes / f64::from(1 << 30),
                    hours / f64::from(24),
                    msg_json["result"]["events"]["buy_storage.buyer"][0]
                        .as_str()
                        .unwrap(),
                    msg_json["result"]["events"]["buy_storage.receiver"][0]
                        .as_str()
                        .unwrap(),
                    msg_json["result"]["events"]["tx.hash"][0]
                        .as_str()
                        .unwrap()
                );
            }
        }
    } else if cmd == 3 {
        return;
    }
}
