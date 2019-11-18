
pub fn send() -> Vec<&'static str>{
    let mut msgs:Vec<&'static str> = vec![];
    msgs.push("USER x@x.com");
    msgs.push("RDEL");
    msgs.push("RCVD");
    msgs.push("QUIT");
    return msgs;
} 
