use crate::control::state::State;

pub struct Help {}

impl super::Reducer for Help {
    fn reduce(&self, mut state: State) -> Result<State, super::ReducerError> {
        Ok(state)
    }
}
