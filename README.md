# .inder

> Chicken Masala legende Wollmilchsau

## order state flow

open -> users can add, modify and delete orders (optional Bestellschluss)
finalized -> order cannot be changed anymore and text is generated, status can
        be queried, total price is available
ordered -> order placed at restaurant (optional ETA)
delivered -> paypal me link of user who paid is posted


## dev

You need to create an empty index dir before starting.

Then load the sangam menu:
```bash
curl -X PUT --data "@sangam.json" -H 'Content-Type: application/json' -v http://localhost:8080/menu/sangam
```

Then try the fuzzy search:
```bash
curl http://localhost:8080/menu/sangam\?search_string\=vindaloo
```

## TODO

- implement status command in matrix
- Prometheus exporter
- Grafana Dashboard
- BIP (Brutto Inder Produkt)
- Deadline as (optional) parameter when starting order
- fuzzy search with search index
- paste paypal.me link of person who posted arrived

