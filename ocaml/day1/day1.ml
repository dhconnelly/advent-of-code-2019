open List
open Printf
open Sys

let read_line ic =
  try Some (input_line ic) with End_of_file -> None

let rec map_lines f ic =
  match read_line ic with
  | Some line -> (f line)::(map_lines f ic)
  | None -> []

let getarg n =
  try argv.(n) with Invalid_argument _ -> failwith "Usage: day1.ml <input_file>"

let fuel mass = mass / 3 - 2

let rec fuelrec mass =
  let x = (fuel mass) in
  if x < 0 then 0 else x + fuelrec x

let process f path = 
  let ic = open_in path in
  let nums = map_lines (fun line -> line |> int_of_string |> f) ic in
  fold_right (+) nums 0

let _ =
  let path = getarg 1 in
  printf "%d\n" (process fuel path);
  printf "%d\n" (process fuelrec path)
