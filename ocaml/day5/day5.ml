open Array
open Printf
include Vm

let getarg n =
  try Sys.argv.(n) with Invalid_argument _ -> failwith "Usage: day2.ml <input_file>"

type opcode = Add | Mul | Halt

let decode = function
  | 1 -> Add
  | 2 -> Mul
  | 99 -> Halt
  | o -> failwith (sprintf "invalid opcode: %d" o)

type state = Running | Halted

type vm = {
  pc: int;
  data: data;
  state: state;
}

let ld pos data = get (get pos data) data
let store pos x data = set x (get pos data) data

let step {pc; data} =
  let (a, b, z) = (ld (pc+1) data), (ld (pc+2) data), (pc+3) in
  match get pc data |> decode with
  | Add -> {pc=pc+4; data=(store z (a + b) data); state=Running}
  | Mul -> {pc=pc+4; data=(store z (a * b) data); state=Running}
  | Halt -> {pc=pc+1; data; state=Halted}

let rec run vm = match vm with
  | {state=Halted} -> vm
  | _ -> run (step vm)




let run_program program noun verb =
  let data = copy program |> set noun 1 |> set verb 2 in
  let vm = run {pc=0; data; state=Running} in
  get 0 (run vm).data

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
