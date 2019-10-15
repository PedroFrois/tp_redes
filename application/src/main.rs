mod smtp;
use smtp::client;
use smtp::server;
//use smtp::server::storage::email::Email;

fn main() {
    let msgs = client::interface::send("teste@sadas.com","teste@teste.com","jfhansdjfuiasdnfuijasdnfsjfui   \n   nasdfuibasdfibasdfuiasdnfsjdfnasdjfuinasdifjpnasdjfuisdabnfyhasdbpnfjasdnfasdujfbaspdij fhdasfu adshmfuasdhfasudmhfasdunfah dufhasdfuasndfoasdnfasdufhasdpufijsdaopifhadpasdf");
    let mut e = server::interface::init();
    for msg in msgs{
        e = server::interface::parse_msg(&msg,e);
    }
    println!("{} ",e.domain);
    println!("{} ",e.to);
    println!("{} ",e.from);
    println!("{} ",e.data);
}
