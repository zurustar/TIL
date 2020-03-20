#!/usr/bin/env python3

def solve(q):
  ns, os = str(q), '+-*/'
  a = []
  for l0 in ['','(','((','(((']:
    for o0 in os:
      for l1 in ['','(','((']:
        for r1 in ['',')']:
          for o1 in os:
            for l2 in ['','(']:
              for r2 in ['',')','))']:
                for o2 in os:
                  for r3 in ['',')','))',')))']:
                    f=l0+ns[0]+o0+l1+ns[1]+r1+o1+l2+ns[2]+r2+o2+ns[3]+r3
                    try:
                      ans=eval(f)
                      if (9.999<ans) and(ans<10.001):
                        return f
                    except:
                      pass
  return None

for i in range(1000,9999):
  f = solve(i)
  if f != None:
    print(f)
