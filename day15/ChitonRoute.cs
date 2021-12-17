using System.Drawing;
using System.Collections.Generic;
using System.Linq;

class ChitonRoute
{
    private readonly int[,] riskLevels;
    private List<Point> points = new List<Point>();

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
                points.Add(new Point(x, y));
            }
        }
    }

    public Point Start { get; set; }
    public Point End { get; set; }

    public int CaveSize { get { return riskLevels.GetLength(0) * riskLevels.GetLength(1); } }

    public int LowestTotalRisk()
    {
        int result = Dijsktra();
        Console.WriteLine("Result is " + result);
        return result;
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

    private int Dijsktra()
    {
        var distances = new PriorityQueue<RiskToPoint, int>();
        var spt = new List<RiskToPoint>();

        // initialize distances to infinite
        foreach (var p in points)
        {
            distances.Enqueue(new RiskToPoint(p, int.MaxValue), int.MaxValue);
        }

        // start distance is 0
        distances.Enqueue(new RiskToPoint(Start, 0), 0);

        // find shortest path of all vertices
        while (spt.Count < CaveSize)
        {
            // pick the minimum distance from the set of vertices not yet processed
                        
            var next = distances.Dequeue();
            if (! spt.Any(rtp => rtp.Point == next.Point))
            {
                // mark the point as visited
                Console.WriteLine($"Locking in {next.Point} at {next.Risk}");
                spt.Add(next);
                // update the distance value of adjacents
                foreach (var neighbor in Adjacents(next.Point).Where(neighbor => InBounds(neighbor)))
                {
                    var totalWeight = next.Risk + Risk(neighbor);
                    distances.Enqueue(new RiskToPoint(neighbor, totalWeight), totalWeight);
                    Console.WriteLine($" scoring {neighbor} at {totalWeight}");
                }
            }
        }

        return spt.Find(rtp => rtp.Point == End).Risk;
    }

    class RiskToPoint {
        public Point Point { get; init; }
        public int Risk { get; set; }

        public RiskToPoint(Point destination, int risk)
        {
            Point = destination;
            Risk = risk;
        }
    }
}  