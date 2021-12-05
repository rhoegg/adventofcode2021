%dw 2.0
import mergeWith from dw::core::Objects
import drop from dw::core::Arrays
output application/json

var lines = day5
    splitBy  '\n'
        map (
            $ splitBy /\s*\->\s*/
            map ($ splitBy ',')
            map {x: $[0], y: $[1]}
        )
var puzzle1Lines = lines 
    filter (
        ($[0].x == $[1].x) or ($[0].y == $[1].y)
    ) 
var puzzle2AddedLines = lines 
    filter ($[0].x != $[1].x) and ($[0].y != $[1].y)


fun pointsInLine(line, index) =
    if (line[0].x == line[1].x or line[0].y == line[1].y)
        pointsInPuzzle1Line(line)
    else
        pointsInPuzzle2Line(line)

fun pointsInPuzzle1Line(line) = 
    if (line[0].x == line[1].x)
        (line[0].y to line[1].y) map (line[0].x ++ " " ++ $)
    else
        (line[0].x to line[1].x) map ($ as String ++ " " ++ line[0].y)

fun pointsInPuzzle2Line(line) =
    (line[0].x to line[1].x) zip (line[0].y to line[1].y)
        map ($[0] as String ++ " " ++ $[1])


fun countOverlaps(lines, puzzleFunc) = sizeOf(
    (lines flatMap puzzleFunc($))
    map ($.x ++ " " ++ $.y)
    groupBy $
    pluck $
    filter sizeOf($) > 1
)

var thermalPoints = 
    lines flatMap (line, index) -> pointsInLine(line, index)

fun countRepeats(points, previousPoints = []) =
    if (isEmpty(points)) 0
    else
        if (previousPoints contains points[0]) 1 else 0
        + countRepeats(points drop 1, previousPoints ++ [points[0]])

---
{
    // p1: countOverlaps(puzzle1Lines, pointsInPuzzle1Line),
    // p2: countOverlaps(puzzle2AddedLines, pointsInPuzzle2Line),
    // attempt: countOverlaps(lines, pointsInLine),
    // points: sizeOf(lines flatMap pointsInLine($)),
    // distinctPoints: sizeOf(lines flatMap pointsInLine($) distinctBy ($)),
    // overlapPoints: sizeOf(thermalPoints) - sizeOf(thermalPoints distinctBy ($))
    // allPoints: sizeOf(thermalPoints),
    // firstCouple: sizeOf(thermalPoints filter ($.index < 3)),
    // unique: sizeOf(thermalPoints distinctBy $)
    // countRepeats: countRepeats(thermalPoints)
    finalAttempt: sizeOf(
        (thermalPoints groupBy $) 
            pluck sizeOf($)
            filter $ > 1
    )
}

