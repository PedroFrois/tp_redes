use super::settings; 

pub fn send(from:&str,to:&str,data:&str) -> Vec<String>{
    let mut msgs:Vec<String> = vec![];
    msgs.push(format!("HELO {}",settings::DOMAIN).to_string());
    msgs.push(format!("MAIL FROM:{}",from).to_string());
    msgs.push(format!("RCPT TO:{}",to).to_string());
    let mut i=0;
    while i < data.len(){
        if i+20 >= data.len(){
            msgs.push(format!("DATA {}",&data[i..]).to_string());
        }else{
            msgs.push(format!("DATA {}",&data[i..i+20]).to_string());
        }
        i+=20;
    }
    msgs.push(format!("QUIT").to_string());
    return msgs;
} 
