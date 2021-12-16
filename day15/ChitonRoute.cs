using System.Drawing;

class ChitonRoute
{
    private readonly int[,] riskLevels;
    private Dictionary<Route, int> riskCache = new Dictionary<Route, int>();
    public ChitonRoute(string[] inputLines)
    {
        var riskLevelInput = new List<List<string>>();

        foreach (var inputLine in inputLines)
        {
            var riskLevelRow = new List<string>();
            foreach (var riskChar in inputLine)
            {
                riskLevelRow.Add(riskChar.ToString());
            }
            riskLevelInput.Add(riskLevelRow);
        }
        Console.WriteLine("Cavern is " + riskLevelInput[0].Count + "x" + riskLevelInput.Count);

        riskLevels = new int[riskLevelInput.Count, riskLevelInput[0].Count];

        for (int y = 0; y < riskLevels.GetLength(0); y++)
        {
            for (int x = 0; x < riskLevels.GetLength(1); x++)
            {
                riskLevels[y,x] = int.Parse(riskLevelInput[y][x]);
            }
        }
    }

    public Point Start { get; set; }
    public Point End { get; set; }

    public int LowestTotalRisk()
    {
        int result =  LowestTotalRisk(new Point(0, 0), new Point(riskLevels.GetLength(1) - 1, riskLevels.GetLength(0) - 1), new List<Point>());
        Console.WriteLine("Result is " + result);
        return result;
    }

    private int LowestTotalRisk(Point start, Point end, List<Point> visited)
    {
        if (start == end)
        {
            return 0;
        }
        if (IsAdjacent(start, end))
        {
            return Risk(end);
        }
        var newVisited = new List<Point>(visited);
        newVisited.Add(start);
        // newVisited.Add(end);

        var neighborRisk = new Dictionary<Point, int>();
        foreach (var neighbor in Adjacents(start))
        {
            if (InBounds(neighbor) && ! visited.Contains(neighbor))
            {
                var route = new Route(neighbor, end);

                // 7,3 never checks 7,4
                if (riskCache.ContainsKey(route))
                {
                    neighborRisk[neighbor] = riskCache[route];
                    newVisited.Add(end);
                    Console.WriteLine("C " + start + " neighbor " + neighbor + " risk is " + Risk(neighbor) + " + " + (neighborRisk[neighbor] - Risk(neighbor)));
                }
                else
                {
                    var lowest = LowestTotalRisk(neighbor, end, newVisited);
                    if (lowest > 0)
                    {
                        Console.WriteLine("- " + start + " neighbor " + neighbor + " risk is " + Risk(neighbor) + " + " + lowest);
                        neighborRisk[neighbor] = Risk(neighbor) + lowest;
                        riskCache[route] = neighborRisk[neighbor];
                    }
                }
            } else { if (InBounds(neighbor)) { Console.WriteLine("V " + start + " neighbor " + neighbor); } }
        }
        neighborRisk = neighborRisk.Where(x => x.Value > 0).ToDictionary(e => e.Key, e => e.Value);
        if (neighborRisk.Count == 0)
        {
            return -1;
        }
        return neighborRisk.Min(x => x.Value);
    }

    private int Risk(Point pos)
    {
        return riskLevels[pos.Y, pos.X];
    }

    private bool InBounds(Point pos)
    {
        if (pos.X < 0 || pos.X >= riskLevels.GetLength(1))
        {
            return false;
        }
        if (pos.Y < 0 || pos.Y >= riskLevels.GetLength(0))
        {
            return false;
        }

        return true;
    }

    private bool IsAdjacent(Point p1, Point p2)
    {
        if (p1.X == p2.X)
        {
            if (Math.Abs(p1.Y - p2.Y) == 1)
            {
                return true;
            }
        }
        if (p1.Y == p2.Y)
        {
            if (Math.Abs(p1.X - p2.X) == 1)
            {
                return true;
            }
        }
        return false;
    }

    private List<Point> Adjacents(Point p)
    {
        var result = new List<Point>();
        foreach (var offset in new int[]{-1, 1})
        {
            result.Add(new Point(p.X, p.Y + offset));
            result.Add(new Point(p.X + offset, p.Y));
        }
        return result;
    }

    public IEnumerable<Waypoint> PathsToEnd()
    {
        return PathsTo(End);
    }
    
    public IEnumerable<Waypoint> PathsTo(Point destination, Waypoint? next = null)
    {
        Waypoint wp = new Waypoint(destination, next);
        Console.WriteLine(wp.Distance + " far away");
        if (destination == Start) return new[] { wp };
        if (wp.Distance > (riskLevels.GetLength(0) * riskLevels.GetLength(1) / 4)) return Enumerable.Empty<Waypoint>();
        
        return Adjacents(destination)
            .Where(InBounds)
            .Where(p => ! wp.LeadsTo(p))
            .SelectMany(p => PathsTo(p, wp));
    }

    public struct Route
    {
        public Route(Point start, Point end)
        {
            Start = start;
            End = end;
        }

        public Point Start { get; init; }
        public Point End { get; init; }

        public override string ToString() => $"{Start} - {End}";
    }

    public class Waypoint 
    {
        public Waypoint(Point here, Waypoint? next)
        {
            Here = here;
            Next = next;
        }

        public Point Here { get; init; }
        
        public Waypoint? Next { get; init; }

        public bool LeadsTo(Point destination)
        {
            if (destination == Here) return true;
            if (null == Next) return false;
            return Next.LeadsTo(destination);
        }

        public int Distance { 
            get {
                if (null == Next) return 0;
                return 1 + Next.Distance;
            }
        }
    }
}