use crate::control::Action;

pub struct Finalize {}

impl Action for Finalize {
    fn reduce(&self,state: &mut crate::control::State) -> Result<(),crate::control::ReducerError> {
        Ok(())
    }
}

