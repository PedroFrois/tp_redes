use super::super::storage::email::Email;
use super::super::storage::control;
use super::super::settings;
use super::msg;

static mut EMAIL:Email = Email{domain:"",to:"",from:"",data:""};

pub fn init(){
    unsafe{
        EMAIL = Email{domain:"",to:"",from:"",data:""};
    }
}

pub fn parse_msg(msg:&'static str) -> String{
    let space_index = msg.find(' ').unwrap();
    let command = &msg[..space_index-1];
    let content = &msg[space_index+1..];
    let resp_code = match command{
        "HELO" => helo(content),
        "MAIL" => mail(content),
        "RCPT" => rcpt(content),
        "DATA" => data(content),
        "RSET" => rset(),
        "QUIT" => quit(),
        "NOOP" => noop(),
        _      => 500,
    };
    return msg::get_resp_msg(resp_code);
}

fn helo(domain:&'static str) -> i32{
    unsafe{
        EMAIL.domain = domain;
    }
    return 220;
}

fn mail(from:&'static str) -> i32{
    let aux:Vec<&str> = from.split(":").collect();
    unsafe{
        EMAIL.from = aux[1];
    }
    return 250;
}

fn rcpt(to:&'static str) -> i32{
    let aux:Vec<&str> = to.split(":").collect();
    unsafe{
        EMAIL.to = aux[1];
        let d:Vec<&str> = EMAIL.to.split('@').collect();
        if  d[1] == settings::DOMAIN {
            return 250;
        }else{
            return 251;
        }
    }
}

fn data(line:&'static str) -> i32{
    unsafe{
        EMAIL.data = line;
    }
    return 250;
}

fn rset() -> i32{
    unsafe{
        EMAIL.from = "";
        EMAIL.to = "";
        EMAIL.data = "";
        EMAIL.domain = "";
    }
    return 250;
}

fn noop() -> i32{
    return 250;
}

fn quit() -> i32{
    unsafe{
        let d:Vec<&str> = EMAIL.to.split('@').collect();
        if  d[1] == settings::DOMAIN {
            store();
        }else{
            forward();
        }
    }
    return 221
}

fn forward(){
    print!("chamando cliente");
}

fn store(){
    print!("save");
    unsafe{
        control::save(&EMAIL);
        rset();
    }
}
