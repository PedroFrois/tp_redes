use super::super::settings;

pub fn get_resp_msg(code:i32) -> String {
    match code{
        500 => "500 Syntax error, command unrecognized".to_string(),
        501 => "501 Syntax error in parameters or arguments".to_string(),
        502 => "502 Command not implemented".to_string(),
        503 => "503 Bad sequence of commands".to_string(),
        504 => "504 Command parameter not implemented".to_string(),
        211 => "211 System status, or system help reply".to_string(),
        214 => "214 Help message".to_string(),
        220 => format!("220 {} Service ready",settings::DOMAIN),
        221 => format!("221 {} Service closing transmission channel",settings::DOMAIN),
        421 => format!("421 {} Service not available, closing transmission channel",settings::DOMAIN),
        250 => "250 Requested mail action okay, completed".to_string(),
        251 => "251 User not local; will forward ".to_string(),
        450 => "450 Requested mail action not taken: mailbox unavailable".to_string(),
        550 => "550 Requested action not taken: mailbox unavailable".to_string(),
        451 => "451 Requested action aborted: error in processing".to_string(),
        551 => "551 User not local; please try to domain".to_string(),
        452 => "452 Requested action not taken: insufficient system storage".to_string(),
        552 => "552 Requested mail action aborted: exceeded storage allocation".to_string(),
        553 => "553 Requested action not taken: mailbox name not allowed".to_string(),
        354 => "354 Start mail input; end with <CRLF>.<CRLF>".to_string(),
        554 => "554 Transaction failed".to_string(),
        _   => "CODE NOT SUPPORTED".to_string(),
    }
}
