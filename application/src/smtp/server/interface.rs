use super::storage::email::Email;
use super::storage::control;
use super::settings;


pub fn init() -> Email{
    return rset();
}

pub fn parse_msg(msg:&str,e:Email) -> Email{
    let s:Vec<&str> = msg.split_whitespace().collect();
    match s[0]{
        "HELO" => helo(s[1].to_string(),e),
        "MAIL" => mail(s[1],e),
        "RCPT" => rcpt(s[1],e),
        "DATA" => data(s[1..].join(" "),e),
        "RSET" => rset(),
        "QUIT" => quit(e),
        _ => noop(e),
    }
}

fn helo(domain:String,mut e:Email) -> Email{
    e.domain = domain;
    return e;
}

fn mail(from:&str,mut e:Email) -> Email{
    let aux:Vec<&str> = from.split(":").collect();
    e.from = aux[1].to_string();
    return e;
}

fn rcpt(to:&str,mut e:Email) -> Email{
    let aux:Vec<&str> = to.split(":").collect();
    e.to = aux[1].to_string();
    return e;
}

fn data(line:String,mut e:Email) -> Email{
    e.data = e.data + &line;
    return e;
}

fn rset() -> Email{

    return Email{from:"".to_string(),to:"".to_string(),data:"".to_string(),domain:"".to_string()};
}

fn noop(e:Email) -> Email{
    return e;
}

fn quit(e:Email) -> Email{
    let d:Vec<&str> = e.to.as_str().split('@').collect();
    if  d[1] == settings::DOMAIN {
        return store(e);
    }else{
        return forward(e);
    }
}

fn forward(e:Email) -> Email{
    print!("chamando cliente");
    return e;
}

fn store(e:Email) -> Email{
    print!("save");
    control::save(e);
    return rset();
}
