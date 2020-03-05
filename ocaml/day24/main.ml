open Printf
open Pt2
module PtMap = Map.Make(Pt2)

type state = Alive | Dead
type grid = state PtMap.t

let tile_of_char : char -> state = function
  | '#' -> Alive
  | '.' -> Dead
  | ch -> failwith (sprintf "invalid state: %c" ch)

let read_grid (ic : in_channel) : grid =
  let read_tile row (col, ch) = (col, row), tile_of_char ch in
  let read_row row : (Pt2.t * state) Seq.t =
    input_line ic |> String.to_seqi |> Seq.map (read_tile row) in
  let rec loop row acc =
    try PtMap.add_seq (read_row row) acc |> loop (row+1)
    with End_of_file -> acc in
  loop 0 PtMap.empty

let print_grid (g : grid) =
  let mc, mr = PtMap.fold (fun (c,r) _ (mc,mr) -> max mc c, max mr r) g (0,0) in
  for row=0 to mr do
    for col=0 to mc do
      printf "%c" (match PtMap.find (col,row) g with
      | Alive -> '#'
      | Dead -> '.')
    done;
    printf "\n"
  done

let () = open_in Sys.argv.(1) |> read_grid |> print_grid
