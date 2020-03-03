open Printf
open Pt2
module PtMap = Map.Make(Pt2)

type tile = Wall | Entrance | Passage | Door of char | Key of char

let key_of ch = Char.lowercase_ascii ch
let door_of ch = Char.uppercase_ascii ch

type maze = {grid: tile PtMap.t; rows: int; cols: int}

let print_tile = function
  | Wall -> '#'
  | Entrance -> '@'
  | Passage -> '.'
  | Door ch -> ch
  | Key ch -> ch

let print m =
  for row=0 to m.rows-1 do
    for col=0 to m.cols-1 do
      print_tile (PtMap.find (col, row) m.grid) |> print_char
    done;
    print_newline ()
  done

let parse_tile row col c =
  let tile = match c with
  | '#' -> Wall
  | '@' -> Entrance
  | '.' -> Passage
  | 'a'..'z' as c -> Key c
  | 'A'..'Z' as c -> Door c
  | c -> failwith (sprintf "bad tile: %c" c) in
  (col, row), tile

let read_chars ic = input_line ic |> String.to_seq |> List.of_seq
let parse_tiles row chars = List.mapi (parse_tile row) chars
let add_point grid (pt, tile) = PtMap.add pt tile grid

let read_maze ic =
  let rec loop row m =
    try
      let pts = read_chars ic |> parse_tiles row in
      let grid = List.fold_left add_point m.grid pts in
      let m = {grid; rows=row+1; cols=List.length pts} in
      loop (row+1) m
    with End_of_file -> m in
  loop 0 {grid=PtMap.empty; rows=0; cols=0}

let () =
  let m = open_in Sys.argv.(1) |> read_maze in
  print m
