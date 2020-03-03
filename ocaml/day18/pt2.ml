type t = int * int

let compare (x1,y1) (x2,y2) =
  match Stdlib.compare x1 x2 with
  | 0 -> Stdlib.compare y1 y2
  | c -> c

let nbrs (x,y) = [(x+1, y); (x-1, y); (x, y+1); (x, y-1)]

let fmt (x,y) = Printf.sprintf "(%d, %d)" x y
