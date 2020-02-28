open Printf

type data = int array

let get pos data = data.(pos)
let set x pos data = data.(pos) <- x; data
let copy = Array.copy

let read ic =
  let next_instr () =
    try Some (Scanf.bscanf ic "%d%c" (fun i _ -> i))
    with End_of_file -> None
  in let rec read_acc acc = match next_instr () with
  | None -> acc
  | Some instr -> read_acc (instr::acc)
  in read_acc [] |> List.rev |> Array.of_list

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

let vm_new data = {pc=0; data; state=Running}
let vm_data vm = vm.data
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
