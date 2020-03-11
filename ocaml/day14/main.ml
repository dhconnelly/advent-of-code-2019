module StrMap = Map.Make(String)
type quant = {chem: string; amt: int}
type reaction = {out: quant; ins: quant list}
type reaction_map = reaction StrMap.t
type chem_map = int StrMap.t

let reaction_delim = Str.regexp_string " => "
let quants_delim = Str.regexp_string ", "

let read_reactions (ic: in_channel): reaction_map =
  let parse_quant s = Scanf.sscanf s "%d %s" (fun amt chem -> {amt; chem}) in
  let parse_quants qs = Str.split quants_delim qs |> List.map parse_quant in
  let parse_reaction line = match Str.split reaction_delim line with
  | ins::out::[] -> {ins=parse_quants ins; out=parse_quant out}
  | _ -> failwith ("bad reaction: " ^ line) in
  let add_reaction m r = StrMap.add r.out.chem r m in
  let rec loop acc =
    try input_line ic |> parse_reaction |> add_reaction acc |> loop
    with End_of_file -> acc in
  loop StrMap.empty

let print_reactions (rs: reaction_map) =
  StrMap.iter (fun _ {out; ins} ->
    Printf.printf "%d %s: " out.amt out.chem;
    List.iter (fun {amt; chem} -> Printf.printf "%d %s, " amt chem) ins;
    print_newline ()) rs

type react_state = {need: chem_map; have: chem_map}

let initial_state (rs: reaction_map): react_state =
  {need=StrMap.empty |> StrMap.add "FUEL" 1; have=StrMap.empty}

let print_state ({need; have}: react_state) =
  let print_chems chem amt = Printf.printf "%s: %d\n" chem amt in
  Printf.printf "Need:\n";
  StrMap.iter print_chems need;
  Printf.printf "Have:\n";
  StrMap.iter print_chems have

let rec reduce (rs: reaction_map) (state: react_state): react_state =
  state

let cost (rs: reaction_map): int =
  0

let () =
  let reactions = Sys.argv.(1) |> open_in |> read_reactions in
  let state = initial_state reactions in
  print_state state
