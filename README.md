# .inder

> Chicken Masala legende Wollmilchsau

## UI pages

- login/logout/signup
- admin area (user management)
- profile page (link matrix account, set username, ...)
- startpage with tutorial, later dashboards
- menus route with list of all menus
    - single menu detail page
- orders overview page with list of all orders (grouped by restaurant and active
  orders)
    - order detail page, table with orders grouped by user with total price


## order state flow

open -> users can add, modify and delete orders (optional Bestellschluss)
finalized -> order cannot be changed anymore and text is generated, status can
        be queried, total price is available
ordered -> order placed at restaurant (optional ETA)
delivered -> paypal me link of user who paid is posted


## Getting Started

open dev shell

```bash
nix flake --impure
```

start database

```bash
devenv up
```

(optional) run database migrations manually

```bash
migrate -database ${DATABASE_URL} -path db/migrations up
```

## TODO

- implement status command in matrix
- Prometheus exporter
- Grafana Dashboard
- BIP (Brutto Inder Produkt)
- Deadline as (optional) parameter when starting order
- fuzzy search with search index
- paste paypal.me link of person who posted arrived
- command for received payments when the order arrived
- only the sugar person can issue "paid" commands
- "kneecap list" with persons who haven't paid their order items
- only the initiator can modify the order

