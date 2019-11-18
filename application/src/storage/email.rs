pub struct Email{
    pub domain: &'static str,
    pub from: &'static str,
    pub to: &'static str,
    pub data: &'static str,
}

impl Email{
    pub fn to_string(&self) -> String{
        return format!("From:{}\nTo:{}\n\n{}\n",self.from,self.to,self.data);
    }
}

