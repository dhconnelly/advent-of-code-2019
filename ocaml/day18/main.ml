open Printf
open Pt2
module PtMap = Map.Make(Pt2)
module PtSet = Set.Make(Pt2)
module CharMap = Map.Make(Char)

(* tiles *)
type tile = Wall | Entrance | Passage | Door of char | Key of char

let parse_tile row col c =
  let tile = match c with
  | '#' -> Wall
  | '@' -> Entrance
  | '.' -> Passage
  | 'a'..'z' as c -> Key c
  | 'A'..'Z' as c -> Door c
  | c -> failwith (sprintf "bad tile: %c" c) in
  (col, row), tile

let print_tile = function
  | Wall -> '#'
  | Entrance -> '@'
  | Passage -> '.'
  | Door ch -> ch
  | Key ch -> ch

let is_passable = function
  | Entrance | Passage -> true
  | _ -> false
let is_door_or_key = function
  | Door _ | Key _ -> true
  | _ -> false
let char_of = function
  | Door ch | Key ch -> ch
  | _ -> failwith "not a char tile"

(* grid *)
type grid = tile PtMap.t
let add_point g (pt, tile) = PtMap.add pt tile g

let read_grid ic =
  let rec loop row g = try
    let chars = input_line ic |> String.to_seq |> List.of_seq in
    let tiles = List.mapi (parse_tile row) chars in
    let g = List.fold_left add_point g tiles in
    loop (row+1) g with End_of_file -> g in
  loop 0 PtMap.empty

let print_grid g =
  let max_col, max_row = PtMap.fold
    (fun (c,r) v (mc,mr) -> max mc c, max mr r) g (0,0) in
  for row=0 to max_row do
    for col=0 to max_col do
      print_tile (PtMap.find (col, row) g) |> print_char
    done;
    print_newline ()
  done

(* bfs *)
type dist_map = int CharMap.t

let print_dists d =
  CharMap.iter (fun ch dist -> printf "%c -> %d\n" ch dist) d

type bfs_node = {pt: Pt2.t; tile: tile; dist: int}

let node_of pt dist g = {pt; tile=PtMap.find pt g; dist}

let bfs_nbrs pt g v dist =
  let nbrs = Pt2.nbrs pt |> List.filter (fun pt -> PtMap.mem pt g) in
  let to_visit = List.filter (fun pt -> PtSet.mem pt v |> not) nbrs in
  let nodes = List.map (fun pt -> node_of pt dist g) to_visit in
  List.filter (fun {tile} -> tile <> Wall) nodes

let mark_visited v {pt} = PtSet.add pt v

let update_dists d {pt; tile; dist} =
  if is_door_or_key tile then CharMap.add (char_of tile) dist d else d

let bfs pt g : dist_map =
  let q = Queue.create () in Queue.add (node_of pt 0 g) q;
  let rec loop v d = match Queue.take_opt q with
    | None -> d
    | Some nd -> visit v d nd
  and visit v d {pt; tile; dist} =
    let nbrs = bfs_nbrs pt g v (dist+1) in
    let v = List.fold_left mark_visited v nbrs in
    let d = List.fold_left update_dists d nbrs in
    List.iter (fun nd -> if is_passable nd.tile then Queue.add nd q) nbrs;
    loop v d in
  loop (PtSet.add pt PtSet.empty) CharMap.empty

(* main *)
let () =
  let g = open_in Sys.argv.(1) |> read_grid in
  print_grid g;
  print_dists (bfs (8,4) g)

