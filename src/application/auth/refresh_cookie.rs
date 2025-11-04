use actix_web::cookie::Cookie;

pub fn create_refresh_cookie(refresh_token: &str) -> Cookie<'static> {
    Cookie::build("refresh_token", refresh_token.to_string())
        .http_only(true)
        .path("/api/auth/refresh")
        .max_age(actix_web::cookie::time::Duration::days(7))
        .same_site(actix_web::cookie::SameSite::Strict)
        .finish()
}

pub fn create_expired_refresh_cookie() -> Cookie<'static> {
    Cookie::build("refresh_token", "")
        .http_only(true)
        .path("/api/auth/refresh")
        .max_age(actix_web::cookie::time::Duration::seconds(0))
        .same_site(actix_web::cookie::SameSite::Strict)
        .finish()
}
