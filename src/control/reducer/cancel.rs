use crate::control::Action;

pub struct Cancel {}

impl Action for Cancel {
    fn reduce(&self,state: &mut crate::control::State) -> Result<(),crate::control::ReducerError> {
        Ok(())
    }
}

