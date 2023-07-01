use std::{collections::HashMap, sync::Arc};

use tokio::sync::{
    mpsc::{self, UnboundedReceiver, UnboundedSender},
    RwLock,
};

use super::{action, state::State};
use crate::control::action::Action;

pub type EffectFn = fn(state: &State);
pub type ActionSender = UnboundedSender<action::ActionEnum>;
pub type ActionReceiver = UnboundedReceiver<action::ActionEnum>;

pub struct Store {
    rx: ActionReceiver,
    tx: ActionSender,
    effects: HashMap<super::action::ActionEnum, Vec<EffectFn>>,
    state: Arc<RwLock<State>>,
}

impl Store {
    pub fn new(state: Arc<RwLock<State>>) -> Self {
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
    pub async fn get_state_snapshot(&self) -> State {
        self.state.read().await.clone()
    }
    pub fn register_effect(&mut self, effect: EffectFn, for_action: super::action::ActionEnum) {
        if let Some(effects) = self.effects.get_mut(&for_action) {
            effects.push(effect);
        } else {
            self.effects.insert(for_action, vec![effect]);
        }
    }
    pub async fn listen(&mut self) {
        loop {
            if let Some(action) = self.rx.recv().await {
                match action.reduce(self.get_state_snapshot().await) {
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
