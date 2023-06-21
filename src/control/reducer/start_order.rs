use crate::control::Action;

pub struct StartOrder {}

impl Action for StartOrder {
    fn reduce(&self,state: &mut crate::control::State) -> Result<(),crate::control::ReducerError> {
        Ok(())
    }
}
