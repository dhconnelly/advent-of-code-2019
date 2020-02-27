open Array
open Printf

let getarg n =
  try Sys.argv.(n) with Invalid_argument _ -> failwith "Usage: day2.ml <input_file>"

let next_instr ic =
  try Some (Scanf.bscanf ic "%d%c" (fun i _ -> i)) with End_of_file -> None

let rec read_program ic =
  match next_instr ic with
  | None -> []
  | Some instr -> instr::(read_program ic)

type opcode = Add | Mul | Halt
type state = Running | Halted

type vm = {
  pc: int;
  data: int array;
  state: state;
}

let get vm pos =
  vm.data.(vm.data.(vm.pc + pos))

let set vm pos x =
  vm.data.(vm.data.(vm.pc + pos)) <- x

let new_vm program = {pc=0; data=of_list program; state=Running}

let decode = function
  | 1 -> Add
  | 2 -> Mul
  | 99 -> Halt
  | o -> failwith (sprintf "invalid opcode: %d" o)

let step vm =
  let {pc; data} = vm in
  if pc > length data then failwith "program terminated without halting"
  else match decode data.(pc) with
  | Add -> set vm 3 ((get vm 1) + (get vm 2)); {vm with pc=pc+4}
  | Mul -> set vm 3 ((get vm 1) * (get vm 2)); {vm with pc=pc+4}
  | Halt -> {vm with state=Halted}

let rec run = function
  | {state=Halted} -> print_endline "halted"
  | vm -> run (step vm)

let run_file path noun verb =
  let ic = Scanf.Scanning.open_in path in
  let vm = read_program ic |> new_vm in
  vm.data.(1) <- noun;
  vm.data.(2) <- verb;
  run vm;
  vm.data.(0)

let () =
  run_file (getarg 1) 12 2 |> printf "%d\n"
