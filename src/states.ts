export enum State {
    IDLE = 0,
    TAKE_ORDERS = 1,
    ORDERED = 2,
    DELIVERED = 3,
}

export enum Transition {
    START_ORDER = 0,
    ADD_ITEM = 1,
    FINALIZE = 2,
    CANCEL = 3,
    ARRIVED = 4
}
