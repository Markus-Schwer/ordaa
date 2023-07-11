use std::{collections::HashMap, sync::Arc};

use tokio::sync::{
    mpsc::{self, UnboundedReceiver, UnboundedSender},
    RwLock,
};

use super::{action::{Action, Reducer}, state::State};

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
                match action.reduce(self.state.read().await.clone()) {
                    Ok(new_state) => {
                        {
                            let mut writable_state = self.state.write().await;
                            *writable_state = new_state;
                        }
                        if let Some(effects) = self.effects.get(&action) {
                            for effect in effects {
                                effect(&self.state.read().await.clone());
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
