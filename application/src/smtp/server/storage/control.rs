use super::email::Email;


pub fn get_all(user:String) -> Vec<Email>{
    print!("{}",user);
    let filtered_emails:Vec<Email> = Vec::new();
    /*for e in emails{
        if e.from == user{
            filtered_emails.push(e);
        }
    }*/
    return filtered_emails;
}

pub fn save(email:Email){
    print!("{}",email.to);
    //emails.push(email);
}
