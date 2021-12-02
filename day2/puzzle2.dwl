%dw 2.0
import countBy, drop, slice from dw::core::Arrays
output application/json

var commands = day2 splitBy '\n'

fun parse(commands) = 
    (commands map (command) -> (command splitBy ' '))
    map (command) -> {
        move: command[0],
        amount: command[1] as Number
    }

fun applyChanges(commands) = parse(commands)
    map (command) -> command match {
        case c if c[0] == "forward" -> { action: "go", amount: c[1] }
        case c if c[0] == "down" -> { action: "aim", amount: c[1] }
        case c if c[0] == "up" -> { action: "aim", amount: c[1] * -1 }
        else -> "nothing"
    }

fun totalDistance(commands) = applyChanges(commands)
    reduce (command, state = {position: [0, 0], aim: 0}) -> {
        position: command.action match {
            case "go" -> [state.position[0] + command.amount, state.position[1] + log(command.amount * state.aim)]
            else -> state.position
        },
        aim: command.action match {
            case "aim" -> state.aim + command.amount
            else -> state.aim
        }
    }

---
totalDistance(commands).position[0] * totalDistance(commands).position[1]

