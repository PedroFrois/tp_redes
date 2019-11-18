#![allow(dead_code)]
mod smtp;
mod storage;
mod settings;
mod pop;

//use smtp::client;
//use smtp::server;
use pop::client as p_client;
use pop::server as p_server;
use storage::email::Email;

fn main() {
    let msgs = p_client::send();
    for msg in msgs{
        let resp = p_server::parse_msg(msg);
        println!("{}",resp);
    }
}
