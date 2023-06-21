use crate::control::Action;

pub struct Help {}

impl Action for Help {
    fn reduce(&self,state: &mut crate::control::State) -> Result<(),crate::control::ReducerError> {
        Ok(())
    }
}


