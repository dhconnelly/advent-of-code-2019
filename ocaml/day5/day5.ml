open Array
open Printf
include Vm

let getarg n =
  try Sys.argv.(n) with Invalid_argument _ -> failwith "Usage: day2.ml <input_file>"

let run_program program noun verb =
  set noun 1 program |> set verb 2 |> vm_new
  |> run |> vm_data |> get 0

let run path =
  let data = Scanf.Scanning.open_in path |> read in
  run_program data 12 2 |> printf "%d\n";
  for noun=0 to 99 do
    for verb=0 to 99 do
      if run_program data noun verb = 19690720 then
        printf "%d\n" (100 * noun + verb)
    done
  done

let () = getarg 1 |> run
