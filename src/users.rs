use std::fmt::Display;

#[derive(Debug, Clone)]
pub struct User {
    name: String
}

impl Display for User {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.name)?;
        Ok(())
    }
}
