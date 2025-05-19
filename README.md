# harmonies
harmonies helper
[rule](https://cdn.svc.asmodee.net/production-libellud/uploads/2024/03/LIBELLUD_HARMONIES_RULES_EN-1.pdf)
## hexagons
[hexmap](https://www.redblobgames.com/grids/hexagons/)

```plain
    s   -r

-q          +q

    +r  -s

assert (s + r + q) == 0
```
so we can use (r, q) to represent a hexagon, s = -q-r if needed.

for any card, we define the `Animal` grid as (0, 0), marking the condition of the card.

## generate
```shell
buf generate pkg/types --template pkg/types/buf.gen.yaml
```
