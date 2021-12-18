using Snailfish;

SnailfishNumber? num = null;
foreach (var line in File.ReadLines("input.txt"))
{
    if (null == num)
    {
        num = SnailfishNumber.Parse(line);
    }
    else
    {
        num = num.Plus(SnailfishNumber.Parse(line));
        num.Reduce();
    }
}

Console.WriteLine("Puzzle 1 sum is " + num.ToString());
Console.WriteLine("Magnitude " + num.Magnitude);

var numbers = File.ReadLines("input.txt").Select(line => SnailfishNumber.Parse(line));

var maxMagnitude = 0;

for (var i = 0; i < (numbers.Count() - 1); i++)
{
    for (var j = i + 1; j < numbers.Count(); j++)
    {
        var sum = numbers.ElementAt(i).Plus(numbers.ElementAt(j));
        sum.Reduce();
        if (sum.Magnitude > maxMagnitude) maxMagnitude = sum.Magnitude;
    }
}

Console.WriteLine("Puzzle 2 magnitude is " + maxMagnitude);
