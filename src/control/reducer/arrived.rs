use crate::control::Action;

pub struct Arrived {}

impl Action for Arrived {
    fn reduce(&self,state: &mut crate::control::State) -> Result<(),crate::control::ReducerError> {
        Ok(())
    }
}
