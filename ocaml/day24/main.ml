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

let iterate nbrs lookup mapi m =
  let is_alive pt = lookup pt m = Some Alive in
  let alive_nbrs pt = nbrs pt |> List.filter is_alive |> List.length in
  let iter pt cur = match cur, alive_nbrs pt with
  | Alive, 1 -> Alive
  | Alive, _ -> Dead
  | Dead, (1 | 2) -> Alive
  | Dead, _ -> Dead in
  mapi iter m

let iterate_flat = iterate Pt2.nbrs PtMap.find_opt PtMap.mapi

let find_repeat (g: grid): int =
  let rec loop g seen =
    let b = pack g in
    if IntSet.mem b seen then b
    else loop {g with m=iterate_flat g.m} (IntSet.add b seen) in
  loop g IntSet.empty

(* assume rows=5 and cols=5 for simplicity in part 2 *)

type rec_pt = {pt: Pt2.t; d: int}
module RecPt = struct
  type t = rec_pt
  let compare {pt=p1; d=d1} {pt=p2; d=d2} =
    if d1 <> d2 then Int.compare d1 d2 else Pt2.compare p1 p2
end
module RecPtMap = Map.Make(RecPt)
type rec_grid = state RecPtMap.t

let rec_grid_of ({m; rows; cols}: grid): rec_grid =
  let m = PtMap.remove (2,2) m in
  let to_rec_pt (pt,st): RecPt.t * state = {pt; d=0}, st in
  PtMap.to_seq m |> Seq.map to_rec_pt |> RecPtMap.of_seq

let rec range lo hi =
  if lo = hi then [] else lo::range (lo+1) hi
let r = range 0 5

let rec_nbrs ({pt=(col, row); d}: RecPt.t): RecPt.t list =
  let left = match col, row with
  | 3, 2 -> List.map (fun row -> {pt=(4, row); d=d+1}) r
  | 0, row -> [{pt=(1, 2); d=d-1}]
  | col, row -> [{pt=(col-1, row); d}] in
  let right = match col, row with
  | 1, 2 -> List.map (fun row -> {pt=(0, row); d=d+1}) r
  | 4, row -> [{pt=(3, 2); d=d-1}]
  | col, row -> [{pt=(col+1, row); d}] in
  let up = match col, row with
  | col, 0 -> [{pt=(2, 1); d=d-1}]
  | 2, 3 -> List.map (fun col -> {pt=(col, 4); d=d+1}) r
  | col, row -> [{pt=(col, row-1); d}] in
  let down = match col, row with
  | 2, 1 -> List.map (fun col -> {pt=(col, 0); d=d+1}) r
  | col, 4 -> [{pt=(2, 3); d=d-1}]
  | col, row -> [{pt=(col, row+1); d}] in
  left @ right @ up @ down

let iterate_rec =
  iterate rec_nbrs RecPtMap.find_opt RecPtMap.mapi

let add_nbrs (rp: rec_pt) _ (g: rec_grid): rec_grid =
  let add_nbr g nbr =
    if RecPtMap.mem nbr g then g else RecPtMap.add nbr Dead g in
  rec_nbrs rp |> List.fold_left add_nbr g

let step g =
  RecPtMap.fold add_nbrs g g |> iterate_rec

let rec steps n g =
  if n = 0 then g else steps (n-1) (step g)

let count_bugs (g: rec_grid): int =
  RecPtMap.fold (fun _ st n -> if st = Alive then n+1 else n) g 0

let () =
  let g = open_in Sys.argv.(1) |> read_grid in
  find_repeat g |> printf "%d\n";
  rec_grid_of g |> steps 200 |> count_bugs |> printf "%d\n"
