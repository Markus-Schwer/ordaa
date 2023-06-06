export enum State {
    IDLE,
    TAKE_ORDERS,
    ORDERED,
    DELIVERED,
}

export enum Transition {
    START_ORDER,
    ADD_ITEM,
    FINALIZE,
    CANCEL,
    ARRIVED
}
