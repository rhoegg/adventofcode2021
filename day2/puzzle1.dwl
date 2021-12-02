%dw 2.0
import countBy, drop, slice from dw::core::Arrays
output application/json

var commands = day2 splitBy '\n'

fun parse(commands) = 
    (commands map (command) -> (command splitBy ' '))
    map (command) -> [command[0], command[1] as Number]

fun distanceChanges(commands) = parse(commands)
    map (command) -> command match {
        case c if c[0] == "forward" -> [c[1], 0]
        case c if c[0] == "down" -> [0, c[1]]
        case c if c[0] == "up" -> [0, c[1] * -1]
        else -> [0, 0]
    }

fun totalDistance(commands) = distanceChanges(commands)
    reduce (command, total = [0, 0]) -> [total[0] + command[0], total[1] + command[1]]

---
totalDistance(commands)[0] * totalDistance(commands)[1]

