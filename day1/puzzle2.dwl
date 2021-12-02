%dw 2.0
import countBy, drop, slice from dw::core::Arrays
output application/json
var measurements = aocInput splitBy '\n' 
    filter (! isEmpty($))
    map $ as Number

var windows = measurements map slice(measurements, $$, $$ + 3) filter sizeOf($) == 3
---
windows zip (windows drop 1) 
    map sum($[0]) < sum($[1])
    countBy $
