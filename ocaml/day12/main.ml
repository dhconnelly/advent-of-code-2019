type vec3 = int * int * int
let zero3 = (0, 0, 0)

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

let pairs (ps: planet list): (planet * planet) list =
  List.combine ps ps |> List.filter (fun (p, q) -> p != q)

let apply_gravity (ps: planet list): planet list =
  let update_vel (gx, gy, gz) (x, y, z) (dx, dy, dz) = 
    (if gx < x then dx-1 else if gx > x then dx+1 else dx),
    (if gy < y then dy-1 else if gy > y then dy+1 else dy),
    (if gz < z then dz-1 else if gz > z then dz+1 else dz) in
  let apply {pos=pos1} {pos=pos2; vel} =
    {pos=pos2; vel=update_vel pos1 pos2 vel} in
  let apply1 ps p =
    List.fold_left (fun acc q -> apply p q::acc) [] ps in
  List.fold_left apply1 ps ps

let apply_velocity (ps: planet list): planet list =
  let apply1 ({pos=(x, y, z); vel=(dx, dy, dz)} as p) =
    {p with pos=(x+dx, y+dy, z+dz)} in
  List.map apply1 ps

let step (ps: planet list): planet list =
  apply_gravity ps |> apply_velocity

let rec steps (n: int) (ps: planet list): planet list =
  if n = 0 then ps else step ps |> steps (n-1)

let energy (ps: planet list): int =
  let energy1 {pos=(x, y, z); vel=(dx, dy, dz)} =
    Int.((abs x + abs y + abs z) * (abs dx + abs dy + abs dz)) in
  List.map energy1 ps |> List.fold_left (+) 0

let () =
  let ps = open_in Sys.argv.(1) |> scan_planets in
  steps 1000 ps |> energy |> Printf.printf "%d\n"
