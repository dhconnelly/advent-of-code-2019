type vec3 = int * int * int
let zero3 = (0, 0, 0)

let print (x, y, z) =
  Printf.sprintf "<x=%3d, y=%3d, z=%3d>" x y z

let scan_vec (line: string): vec3 =
  let to_vec x y z = x, y, z in
  Scanf.sscanf line "<x=%d, y=%d, z=%d>" to_vec

type planet = {pos: vec3; vel: vec3}

let scan_planets (ic: in_channel): planet list =
  let add vs v = {pos=v; vel=zero3}::vs in
  let rec loop vs =
    try input_line ic |> scan_vec |> add vs |> loop
    with End_of_file -> vs in
  loop [] |> List.rev

let print {pos; vel} =
  Printf.printf "pos=%s, vel=%s\n" (print pos) (print vel)

let print_planets (ps: planet list) =
  List.iter print ps

let pairs (ps: planet list): (planet * planet) list =
  List.combine ps ps |> List.filter (fun (p, q) -> p != q)

let apply_gravity (ps: planet list): planet list =
  ps

let apply_velocity (ps: planet list): planet list =
  let apply1 ({pos=(x, y, z); vel=(dx, dy, dz)} as p) =
    {p with pos=(x+dx, y+dy, z+dz)} in
  List.map apply1 ps

let step (ps: planet list): planet list =
  apply_gravity ps |> apply_velocity

let rec steps (n: int) (ps: planet list): planet list =
  if n = 0 then ps else step ps |> steps (n-1)

let () =
  let ps = open_in Sys.argv.(1) |> scan_planets in
  steps 1000 ps |> print_planets
