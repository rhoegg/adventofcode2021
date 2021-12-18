using System.Collections.Generic;

namespace Snailfish
{
    public interface Element
    {
        bool Pair { get; }
        SnailfishNumber? Parent { get; set; }
        SnailfishNumber? FirstPairAtDepth(int depth);
        void AddFromRight(int value);
        void AddFromLeft(int value);
        bool Split();

        int Magnitude { get; }

        string ToString();
    }

    public class RegularNumber : Element
    {
        public RegularNumber(int value)
        {
            Value = value;
        }

        public int Value { get; set; }
        public int Magnitude { get { return Value; } }
        public SnailfishNumber? Parent { get; set; }
        public bool Pair { get { return false; } }
        public SnailfishNumber? FirstPairAtDepth(int depth)
        {
            return null;
        }
        public void AddFromRight(int value)
        {
            Value += value;
        }

        public void AddFromLeft(int value)
        {
            Value += value;
        }

        public bool Split()
        {
            if (Value < 10) return false;
            var pair = new SnailfishNumber();
            pair.Left = new RegularNumber(Value / 2);
            pair.Left.Parent = pair;
            pair.Right = new RegularNumber(Value / 2 + Value % 2);
            pair.Right.Parent = pair;
            pair.Parent = Parent;
            if (this == Parent.Left)
            {
                Parent.Left = pair;
            }
            else
            {
                Parent.Right = pair;
            }
            return true;
        }

        public string ToString()
        {
            return Value.ToString();
        }
    }
    public class SnailfishNumber : Element
    {
        public Element Left { get; set; }
        public Element Right { get; set; }
        public SnailfishNumber? Parent { get; set; }

        public string ToString()
        {
            return $"[{Left.ToString()},{Right.ToString()}]";
        }

        public bool Pair { get { return true; } }

        public int Magnitude 
        { 
            get 
            { 
                return 3 * Left.Magnitude + 2 * Right.Magnitude;
            } 
        }

        public void Reduce()
        {
            var shouldContinue = true;
            while (shouldContinue)
            {
                shouldContinue = Explode() || Split();
            }
        }

        public bool Explode()
        {
            var pair = FirstPairAtDepth(4);
            if (null == pair) return false;

            var l = pair.Left as RegularNumber;
            var r = pair.Right as RegularNumber;
            
            Element current = l;
            while (current.Parent != null && current == current.Parent.Left)
            {
                current = current.Parent;
            }
            if (null != current.Parent)
            {
                current.Parent.Left.AddFromRight(l.Value);
            }
            current = r;
            while (current.Parent != null && current == current.Parent.Right)
            {
                current = current.Parent;
            }
            if (null != current.Parent)
            {
                current.Parent.Right.AddFromLeft(r.Value);
            }
            var zero = new RegularNumber(0);
            zero.Parent = pair.Parent;
            if (pair == pair.Parent.Left)
            {
                pair.Parent.Left = zero;
            }
            else
            {
                pair.Parent.Right = zero;
            }
            return true;
        }

        public SnailfishNumber? FirstPairAtDepth(int depth)
        {
            if (depth == 0) return this;
            foreach (Element next in new Element[] {Left, Right})
            {
                var n = next.FirstPairAtDepth(depth - 1);
                if (null != n) return n;
            }
            return null;
        }

        public void AddFromRight(int value)
        {
            Right.AddFromRight(value);
        }

        public void AddFromLeft(int value)
        {
            Left.AddFromLeft(value);
        }

        public bool Split()
        {
            if (Left.Split()) return true;
            return Right.Split();
        }

        public SnailfishNumber Plus(SnailfishNumber rhs)
        {
            var sum = new SnailfishNumber();
            this.Parent = sum;
            rhs.Parent = sum;
            sum.Left = this;
            sum.Right = rhs;
            return sum;
        }

        public static SnailfishNumber Parse(string notation)
        {
            if (notation[0] != '[') throw new Exception($"panic {notation}");

            var current = new SnailfishNumber();
            var currentRegular = 0;
            
            for (int i = 1; i < notation.Length; i++)
            {
                switch (notation[i])
                {
                    case '[': // new pair
                        var next = new SnailfishNumber();
                        next.Parent = current;
                        current = next;
                        break;
                    case ']': // finished pair
                        if (notation[i - 1] == ']')
                        {
                            current.Parent.Right = current;
                            current = current.Parent;
                        }
                        else
                        {
                            current.Right = new RegularNumber(currentRegular);
                            current.Right.Parent = current;
                            currentRegular = 0;
                        }
                        break;
                    case ',':
                        if (notation[i - 1] == ']')
                        {
                            current.Parent.Left = current;
                            current = current.Parent;
                        }
                        else
                        {
                            current.Left = new RegularNumber(currentRegular);
                            current.Left.Parent = current;
                            currentRegular = 0;
                        }
                        break;
                    default:
                        // regular number
                        currentRegular = 10 * currentRegular + int.Parse(notation[i].ToString());
                        break;
                }                
            }
            return current;
        }
    }
}