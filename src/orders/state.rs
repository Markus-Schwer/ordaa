use std::fmt::Display;

pub enum State {
    // orders of users are collected
    CollectingOrderItems,
    // the orders of the users are communicated to provider
    PlacingOrder,
    // order is in delivery
    InDelivery,
    // order has arrived
    Arrived
}

impl TryFrom<&str> for State {
    type Error = String;
    fn try_from(value: &str) -> Result<Self, Self::Error> {
        match value {
            "CollectingOrders" => Ok(Self::CollectingOrderItems),
            "PlacingOrders" => Ok(Self::PlacingOrder),
            "InDelivery" => Ok(Self::InDelivery),
            "Arrived" => Ok(Self::Arrived),
            _ => Err(format!("'{}' is not a state", value))
        }
    }
}

impl Display for State {
     fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
         write!(f, "{}", match self {
             State::CollectingOrderItems => "CollectingOrders",
             State::PlacingOrder => "PlacingOrders",
             State::InDelivery => "InDelivery",
             State::Arrived => "Arrived",
         })?;
         Ok(())
     }
}
