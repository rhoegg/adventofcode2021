using System.Drawing;

namespace day16
{
    class Bounds
    {
        public Bounds(int l, int r, int b, int t)
        {
            Left = l;
            Right = r;
            Bottom = b;
            Top = t;
        }
        public int Left { get; init; }
        public int Right { get; init; }
        public int Top { get; init; }
        public int Bottom { get; init; }

        public bool Contains(Point p)
        {
            if (Bottom > p.Y) return false;
            if (Top < p.Y) return false;
            if (Left > p.X) return false;
            if (Right < p.X) return false;
            return true;
        }
    }

    class Program
    {
        static void Main(string[] args)
        {
            //x=20..30, y=-10..-5
            // x=241..275, y=-75..-49
            var bounds = new Bounds(241, 275, -75, -49);

            // -lowY to lowY - 1

            var lowx = bounds.Left;
            var dx = 1;

            while  (lowx > 0)
            {
                lowx -= dx++;
            }

            var minVelocityX = dx - 1;
            var maxVelocityX = bounds.Right;
            var minVelocityY = bounds.Bottom;
            var maxVelocityY = -1 * bounds.Bottom - 1;

            bool hitsRectangle(int dx, int dy)
            {
                var pos = new Point(0, 0);
                while (pos.X <= bounds.Right && pos.Y >= bounds.Bottom)
                {
                    pos = new Point( pos.X + dx, pos.Y + dy);
                    if (bounds.Contains(pos)) return true;
                    if (dx > 0) dx -= 1;
                    dy -= 1;
                }
                return false;
            }

            var hits = 0;
            for (var vx = minVelocityX; vx <= maxVelocityX; vx++)
            {
                for (var vy = minVelocityY; vy <= maxVelocityY; vy++)
                {
                    if (hitsRectangle(vx, vy))
                    {
                        hits++;
                    }
                }
            }
            Console.WriteLine(hits);

        }
    }
}
