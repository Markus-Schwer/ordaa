use crate::control::Action;

pub struct Help {}

impl Action for Help {
    fn reduce(
        &self,
        mut state: crate::control::State,
    ) -> Result<crate::control::State, crate::control::ReducerError> {
        Ok(state)
    }
}
