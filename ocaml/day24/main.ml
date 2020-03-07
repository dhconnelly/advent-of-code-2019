open Printf
open Pt2
module PtMap = Map.Make(Pt2)
module IntSet = Set.Make(Int)

type state = Alive | Dead
type grid = {m: state PtMap.t; rows: int; cols: int}

let tile_of_char: char -> state = function
  | '#' -> Alive
  | '.' -> Dead
  | ch -> failwith (sprintf "invalid state: %c" ch)

let read_grid (ic: in_channel): grid =
  let read_tile row (col, ch) = (col, row), tile_of_char ch in
  let read_row row =
    input_line ic |> String.to_seqi |> Seq.map (read_tile row) in
  let rec loop row acc =
    try PtMap.add_seq (read_row row) acc |> loop (row+1)
    with End_of_file -> acc in
  let m = loop 0 PtMap.empty in
  let mc, mr = PtMap.fold (fun (c,r) _ (mc,mr) -> max mc c, max mr r) m (0,0) in
  {m; rows=mr+1; cols=mc+1}

let print_grid (g: grid) =
  for row=0 to g.rows-1 do
    for col=0 to g.cols-1 do
      printf "%c" (match PtMap.find (col,row) g.m with
      | Alive -> '#'
      | Dead -> '.')
    done;
    printf "\n"
  done

let pack (g: grid): int =
  let to_bit (c,r) = function
    | Alive -> Int.shift_left 1 (g.cols*r + c)
    | Dead -> 0 in
  let add_bit pt st b = to_bit pt st |> Int.logor b in
  PtMap.fold add_bit g.m 0

let unpack (x: int) (rows: int) (cols: int): grid =
  let tile row col =
    if Int.(shift_left 1 (row*cols + col) |> logand x) > 0
    then Alive else Dead in
  let rec loop row col m =
    if row = rows then m
    else
      let m = PtMap.add (col,row) (tile row col) m in
      let row = if col+1 = cols then row+1 else row in
      let col = if col+1 = cols then 0 else col+1 in
      loop row col m in
  {m=(loop 0 0 PtMap.empty); rows; cols}

let iterate (g: grid): grid =
  let is_alive pt = PtMap.find_opt pt g.m = Some Alive in
  let alive_nbrs pt = Pt2.nbrs pt |> List.filter is_alive |> List.length in
  let iter pt cur = match cur, alive_nbrs pt with
    | Alive, 1 -> Alive
    | Alive, _ -> Dead
    | Dead, (1 | 2) -> Alive
    | Dead, _ -> Dead in
  {g with m=(PtMap.mapi iter g.m)}

let find_repeat (g: grid): int =
  let rec loop g seen =
    let b = pack g in
    if IntSet.mem b seen then b
    else loop (iterate g) (IntSet.add b seen) in
  loop g IntSet.empty

module RecPt = struct
  type t = {pt: Pt2.t; d: int}
  let compare {pt=p1; d=d1} {pt=p2; d=d2} =
    if d1 <> d2 then Int.compare d1 d2 else Pt2.compare p1 p2
  let fmt {pt; d} = sprintf "(%s at %d)" (Pt2.fmt pt) d
end

(* assume rows=5 and cols=5 for simplicity in part 2 *)

module RecPtMap = Map.Make(RecPt)
type rec_grid = state RecPtMap.t

let rec_grid_of ({m; rows; cols}: grid): rec_grid =
  let to_rec_pt (pt,st): RecPt.t * state = {pt; d=0}, st in
  PtMap.to_seq m |> Seq.map to_rec_pt |> RecPtMap.of_seq

let rec_nbrs ({pt=(col, row); d}: RecPt.t): RecPt.t list =
  let left = match col, row with
  | 3, 2 -> []
  | 0, row -> []
  | col, row -> [] in
  let right = match col, row with
  | 1, 2 -> []
  | 4, row -> []
  | col, row -> [] in
  let up = match col, row with
  | col, 0 -> []
  | 2, 3 -> []
  | col, row -> [] in
  let down = match col, row with
  | 2, 1 -> []
  | col, 4 -> []
  | col, row -> [] in
  left @ right @ up @ down

let print_rec_nbrs g rp =
  rec_nbrs rp |> List.iter (fun rp -> RecPt.fmt rp |> printf "%s\n")

let () =
  let g = open_in Sys.argv.(1) |> read_grid in
  find_repeat g |> printf "%d\n";
  let g = rec_grid_of g in
  print_rec_nbrs g {pt=(3, 3); d=1};
  print_rec_nbrs g {pt=(1, 1); d=0};
  print_rec_nbrs g {pt=(3, 0); d=0};
  print_rec_nbrs g {pt=(4, 0); d=0};
  print_rec_nbrs g {pt=(3, 2); d=1};
  print_rec_nbrs g {pt=(3, 2); d=0}
