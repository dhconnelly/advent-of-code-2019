open Array
open Printf
include Vm

let getarg n =
  try Sys.argv.(n) with Invalid_argument _ -> failwith "Usage: day2.ml <input_file>"

let rec run_vm arg out vm =
  match vm_state vm with
| Running -> run vm |> run_vm out arg
| Input -> run_vm arg out (vm_write arg vm |> run)
| Output -> run_vm arg (vm_read vm) (run vm)
| Halted -> out

let run_program program arg =
  vm_new program |> run_vm 0 arg

let run path =
  let data = Scanf.Scanning.open_in path |> read in
  run_program data 1 |> printf "%d\n";
  run_program data 5 |> printf "%d\n" 

let () = getarg 1 |> run |> ignore
