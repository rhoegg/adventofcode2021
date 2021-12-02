%dw 2.0
output application/json
var measurements = aocInput splitBy '\n' 
    filter (! isEmpty($))
    map $ as Number

fun solve(measurements) = measurements reduce (m, result = {prev: null, count: 0}) -> {
    prev: m, 
    count: result.count + if (m > (result.prev default m)) 1 else 0
} 
---
solve(measurements).count
