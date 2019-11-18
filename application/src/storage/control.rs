use super::email::Email;

pub fn get_all(user:String) -> Vec<&'static Email>{
    print!("{}",user);
    let mut filtered_emails:Vec<&'static Email> = Vec::new();
    //for e in emails{
    //    if e.from == user{
    //        filtered_emails.push(e);
    //    }
    //}
    return filtered_emails;
}
pub fn get_last(user:String) -> Email{
    print!("{}",user);
    let aux:Email = Email{domain:"",to:"",from:"",data:""};
    //for e in emails{
    //    if e.from == user{
    //       aux = e;
    //    }
    //}
    return aux;
}

pub fn save(email:&Email){
    print!("{}",email.to);
    //emails.push(email);
}

pub fn del(email:&Email){
    print!("{}",email.data);
    //emails.delete(email)
}
