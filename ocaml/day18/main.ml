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
let key_of = Char.lowercase_ascii
let door_of = Char.uppercase_ascii

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

type bfs_node = {pt: Pt2.t; tile: tile; dist: int}

let bfs_nbrs pt g v dist =
  let nbrs = Pt2.nbrs pt |> List.filter (fun pt -> PtMap.mem pt g) in
  let to_visit = List.filter (fun pt -> PtSet.mem pt v |> not) nbrs in
  let nodes = List.map (fun pt -> {pt; tile=PtMap.find pt g; dist}) to_visit in
  List.filter (fun {tile} -> tile <> Wall) nodes

let enq q nbr = if is_passable nbr.tile then Queue.add nbr q
let mark_visited v {pt} = PtSet.add pt v
let update_dists d {pt; tile; dist} =
  if is_door_or_key tile then CharMap.add (char_of tile) dist d else d

let bfs pt g =
  let q = Queue.create () in Queue.add {pt; tile=PtMap.find pt g; dist=0} q;
  let rec loop v d = match Queue.take_opt q with
    | None -> d
    | Some nd -> visit v d nd
  and visit v d {pt; tile; dist} =
    let nbrs = bfs_nbrs pt g v (dist+1) in
    let v = List.fold_left mark_visited v nbrs in
    let d = List.fold_left update_dists d nbrs in
    List.iter (enq q) nbrs;
    loop v d in
  loop (PtSet.add pt PtSet.empty) CharMap.empty

let dists g (pt, tile) = match tile with
  | Key ch | Door ch -> Some (ch, bfs pt g)
  | Entrance -> Some ('@', bfs pt g)
  | _ -> None

let all_dists g =
  let pairs = PtMap.to_seq g |> List.of_seq in
  let all_dists = List.filter_map (dists g) pairs in
  List.to_seq all_dists |> CharMap.of_seq

let print_dists d =
  CharMap.iter (fun ch dist -> printf "%c -> %d\n" ch dist) d

let print_all_dists d =
  CharMap.iter (fun ch d -> printf "from %c:\n" ch; print_dists d) d

(* updating dist map *)

let prod xs ys =
  let pair ys x = List.fold_left (fun zs y -> (x,y)::zs) [] ys in
  List.map (pair ys) xs |> List.concat

let connect d ((chl, dl), (chr, dr)) =
  if chl = chr then d
  else if CharMap.mem chr (CharMap.find chl d) then d
  else (
    let dlr = dl + dr in
    let froml, fromr = CharMap.find chl d, CharMap.find chr d in
    let froml, fromr = CharMap.add chr dlr froml, CharMap.add chl dlr fromr in
    CharMap.add chl froml d |> CharMap.add chr fromr)

let connect_across ch d =
  if not (CharMap.mem ch d) then d else
  let nbrs = CharMap.find ch d |> CharMap.to_seq |> List.of_seq in
  let pairs = prod nbrs nbrs in
  List.fold_left connect d pairs

let remove ch d =
  CharMap.remove ch d |> CharMap.map (fun d -> CharMap.remove ch d)

let take_key key d =
  let door = door_of key in
  let d = connect_across key d |> connect_across door in
  let nbrs = CharMap.find key d in
  nbrs, (remove key d |> remove door)

(* bit vector for keys *)

let to_index key = Char.code key - Char.code 'a'
let from_index ix = ix + Char.code 'a' |> Char.chr

let to_bitset keys =
  let codes = List.map to_index keys in
  List.fold_left (fun b n -> Int.logor b (Int.shift_left 1 n)) 0 codes

let bitset_remove b key =
  to_index key |> Int.shift_left 1 |> Int.lognot |> Int.logand b

let of_bitset b =
  let rec loop acc n =
    if n < 0 then acc
    else if Int.logand (Int.shift_left 1 n) b > 0 then loop ((from_index n)::acc) (n-1)
    else loop acc (n-1) in
  loop [] (to_index 'z')

(* shortest path *)

let is_key_nbr (ch, _) = ch >= 'a' && ch <= 'z'

type memo_key = {from: char; need: int}

let rec shortest_path from need d memo =
  let mk = {from; need} in
  match Hashtbl.find_opt memo mk with
  | Some dist -> dist
  | None ->
      let nbrs, d = take_key from d in
      let nbrs = CharMap.to_seq nbrs |> List.of_seq in
      let nbr_keys = List.filter is_key_nbr nbrs in
      if List.length nbr_keys = 0 then 0
      else (
        let dists = List.map (shortest_path_taking need d memo) nbr_keys in
        let min_dist = List.fold_left min Int.max_int dists in
        Hashtbl.replace memo mk min_dist;
        min_dist)
and shortest_path_taking need d memo (ch, dist) =
  let new_need = bitset_remove need ch in
  dist + (shortest_path ch new_need d memo)

(* main *)

let keys m =
  CharMap.to_seq m |> List.of_seq |> List.filter is_key_nbr |>
  List.map (fun (ch, _) -> ch)

let () =
  let g = open_in Sys.argv.(1) |> read_grid in
  let d = all_dists g in
  let memo = Hashtbl.create 26 in
  printf "%d\n" (shortest_path '@' (keys d |> to_bitset) d memo)

