using System;
using System.IO;
using System.Collections.Generic;
using System.Linq;

namespace Day14
{
    class Program
    {
        static void Main(string[] args)
        {
            Polymerization polymerization;

            using (var reader = new StreamReader("input.txt")) 
            {
                var polymerTemplate = reader.ReadLine();
                Console.WriteLine("Polymer Template: " + polymerTemplate);
                polymerization = new Polymerization(polymerTemplate);

                string s = "";
                while ((s = reader.ReadLine()) != null)
                {
                    if (s.Length > 0)
                    {
                        var injectionSpec = s.Split(" -> ");
                        polymerization.AddInsertion(injectionSpec[0], injectionSpec[1]);
                    }
                }
                reader.Close();
            }

            var counts = polymerization.CountChars(40);
            PrintCounts(counts);
            var max = counts.Select(e => e.Value).Max();
            var min = counts.Select(e => e.Value).Min();
            Console.WriteLine("Puzzle 2: " +  (max - min));
        }

        public static void PrintCounts(Dictionary<string, long> counts)
        {
            foreach (var entry in counts)
            {
                Console.WriteLine(entry.Key + " " + entry.Value);
            }
        }
    }

    class Polymerization {
        readonly string polymerTemplate;
        Dictionary<string, string> insertions = new Dictionary<string, string>();
        Dictionary<int, Dictionary<string, Dictionary<string, long>>> cache = 
            new Dictionary<int, Dictionary<string, Dictionary<string, long>>>();
        
        public Polymerization(string template) 
        {
            this.polymerTemplate = template;
        }

        public void AddInsertion(string template, string insertion)
        {
            this.insertions[template] = insertion;
        }
        
        public void Remember(string template, int step, Dictionary<string, long> counts)
        {
            GetStepCache(step)[template] = counts;
        }

        public Dictionary<string, long> GetCached(string template, int step)
        {
            if (! GetStepCache(step).ContainsKey(template))
            {
                return null;
            }
            return GetStepCache(step)[template];
        }

        public string Translate(string pair)
        {
            var insertion = "";
            if (insertions.ContainsKey(pair))
            {
                insertion = insertions[pair];
            }
            return pair[0] + insertion + pair[1];
        }

        private Dictionary<string, Dictionary<string, long>> GetStepCache(int step)
        {
            if (! cache.ContainsKey(step))
            {
                cache[step] = new Dictionary<string, Dictionary<string, long>>();
            }
            return cache[step];
        }

        public Dictionary<string, long> CountChars(int steps)
        {
            return this.CountChars(polymerTemplate, steps);
        }
        public Dictionary<string, long> CountChars(string template, int steps)
        {
            // first letter then expand each pair after
            if (0 == steps) return CountChars(template);

            var pairs = new List<string>();
            for (int i = 1; i < template.Length; i++)
            {
                pairs.Add(template.Substring(i - 1, 2));
            }

            var counts = new Dictionary<string, long>();
            foreach (var pair in pairs)
            {
                var pairCounts = GetCached(pair, steps);
                if (null == pairCounts)
                {
                    pairCounts = CountChars(Translate(pair), steps - 1);
                    Remember(pair, steps, pairCounts);
                    // Console.WriteLine("** " + steps + " " + pair);
                    // Program.PrintCounts(pairCounts);
                }
                    
                counts = SumCounts(counts, pairCounts);
                counts[pair[1].ToString()] -= 1; // don't count the last character because it's the first character of the next
            }
            counts[template[template.Length-1].ToString()] += 1; // put the very last back in
            return counts;
        }

        private Dictionary<string, long> CountChars(string input)
        {
            var counts = new Dictionary<string, long>();
            foreach(var c in input)
            {
                if (counts.ContainsKey(c.ToString()))
                {
                    counts[c.ToString()] += 1;
                }
                else
                {
                    counts[c.ToString()] = 1;
                }
            }
            return counts;
        }

        private Dictionary<string, long> SumCounts(Dictionary<string, long> d1, Dictionary<string, long> d2)
        {
            foreach (var count in d2)
            {
                if (d1.ContainsKey(count.Key))
                {
                    d1[count.Key] += count.Value;
                }
                else
                {
                    d1[count.Key] = count.Value;
                }
            }

            return d1;
        }
    }
}
