use rocket::http::Status;

#[get("/")]
pub fn index() -> Status {
    Status::BadRequest
}
