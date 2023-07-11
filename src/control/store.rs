use std::{collections::HashMap, sync::Arc};

use tokio::sync::{
    mpsc::{self, UnboundedReceiver, UnboundedSender},
    RwLock,
};

use super::{
    action::{Action, Reducer},
    state::State,
};

pub type EffectFn = fn(state: &State);
pub type ActionSender = UnboundedSender<Action>;
pub type ActionReceiver = UnboundedReceiver<Action>;
pub type SharableState = Arc<RwLock<State>>;

pub struct Store {
    rx: ActionReceiver,
    tx: ActionSender,
    effects: HashMap<Action, Vec<EffectFn>>,
    state: SharableState,
}

impl Store {
    pub fn new(state: SharableState) -> Self {
        let (tx, rx) = mpsc::unbounded_channel();
        Self {
            rx,
            tx,
            effects: HashMap::new(),
            state,
        }
    }
    pub fn get_sender(&self) -> ActionSender {
        self.tx.clone()
    }
    pub fn register_effect(&mut self, effect: EffectFn, for_action: Action) {
        if let Some(effects) = self.effects.get_mut(&for_action) {
            effects.push(effect);
        } else {
            self.effects.insert(for_action, vec![effect]);
        }
    }
    pub async fn listen(&mut self) {
        loop {
            if let Some(action) = self.rx.recv().await {
                let mut writable_state = self.state.write().await;
                match action.reduce(writable_state.clone()) {
                    Ok(new_state) => {
                        *writable_state = new_state;
                        if let Some(effects) = self.effects.get(&action) {
                            for effect in effects {
                                effect(&writable_state.clone());
                            }
                        }
                    }
                    Err(err) => panic!("{:?}", err),
                }
            } else {
                panic!("received empty action signal");
            }
        }
    }
}

#[cfg(test)]
mod test {
    use core::panic;
    use std::{matches, time::Duration};
    use tokio::time::{sleep, timeout};

    use crate::control::action::start_order::StartOrder;

    use super::*;

    #[tokio::test]
    async fn should_do_state_transition() {
        let state = Arc::new(RwLock::new(State::new()));
        let mut store = Store::new(state.clone());
        let sender = store.get_sender();
        tokio::spawn(async move {
            store.listen().await;
        });
        if let Err(err) = sender.send(Action::from(StartOrder {})) {
            panic!("could not send message: {}", err);
        }
        loop {
            if matches!(
                timeout(Duration::from_secs(2), state.read())
                    .await
                    .unwrap()
                    .machine_state,
                crate::control::state::MachineState::TakeOrders
            ) {
                return;
            } else {
                sleep(Duration::from_millis(100)).await;
            }
        }
    }
}
