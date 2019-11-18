use super::super::storage::control;
use super::super::storage::email::Email;

static mut SESSION_USER:&str = "";
static mut READ_EMAIL:Email = Email{domain:"",to:"",from:"",data:""};

pub fn parse_msg<'a>(msg:&'static str) -> String{
    let space_index = msg.find(' ').unwrap_or(msg.chars().count()-1);
    let command = &msg[..space_index-1];
    let content = &msg[space_index+1..];
    let resp_msg = match command{
        "USER" => user(content), 
        "RDEL" => rdel(),
        "RCVD" => rcvd(),
        "QUIT" => quit(),
        "RSET" => rset(),
        _      => "ERR".to_string(),
    };
    return resp_msg;
}

fn user(user:&'static str) -> String{
    unsafe{
        SESSION_USER = user.clone();
    }
    return "OK".to_string();
} 

fn rdel() -> String{
    unsafe{
        READ_EMAIL=control::get_last(SESSION_USER.to_string());
        return READ_EMAIL.to_string();
    }
}

fn rcvd() -> String{
    unsafe{
        control::del(&READ_EMAIL)
    }
    return "OK".to_string();
}

fn quit() -> String{
    unsafe{
        SESSION_USER = "";
    }
    return rset();
    
}

fn rset() -> String{
    unsafe{
        READ_EMAIL= Email{domain:"",to:"",from:"",data:""};
    }
    return "OK".to_string();
}
