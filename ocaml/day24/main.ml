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

let () =
  let g = open_in Sys.argv.(1) |> read_grid in
  find_repeat g |> printf "%d\n"
