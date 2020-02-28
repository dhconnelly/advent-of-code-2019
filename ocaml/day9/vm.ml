module IntMap = Map.Make(Int)
type data = int IntMap.t
let empty = IntMap.empty
let set x pos data = IntMap.add pos x data
let get pos data = match IntMap.find_opt pos data with
| None -> 0
| Some x -> x

let read ic =
  let next_instr () =
    try Some (Scanf.bscanf ic "%d%c" (fun i _ -> i)) with End_of_file -> None
  in let rec read_rec () = match next_instr () with
  | None -> []
  | Some instr -> instr::read_rec ()
  in read_rec () |> List.mapi (fun i x -> i,x) |> List.to_seq |> IntMap.of_seq

type opcode =
  Add | Mul | Read | Write | JmpIf | JmpNot | Lt | Eq | Halt | AdjRel
type mode = Pos | Imm | Rel
type instruction = { op: opcode; modes: mode*mode*mode }

let mode_of = function
  | 0 -> Pos
  | 1 -> Imm
  | 2 -> Rel
  | n -> failwith (Printf.sprintf "invalid mode: %d" n)

let modes_of x =
  let x = x/100 in
  let (a,b,c) = (x mod 10), (x/10 mod 10), (x/100 mod 10) in
  (mode_of a), (mode_of b), (mode_of c)

let decode x =
  let op = match x mod 100 with
  | 1  -> Add
  | 2  -> Mul
  | 3  -> Read
  | 4  -> Write
  | 5  -> JmpIf
  | 6  -> JmpNot
  | 7  -> Lt
  | 8  -> Eq
  | 9  -> AdjRel
  | 99 -> Halt
  | o -> failwith (Printf.sprintf "invalid opcode: %d" o)
  in {op; modes=modes_of x}

type state = Running | Halted | Input | Output

type vm = {
  pc: int;
  data: data;
  state: state;
  input: int;
  output: int;
  rel: int;
}

let vm_empty = {pc=0; data=empty; input=0; output=0; state=Halted; rel=0}
let vm_new data = {vm_empty with data=data; state=Running}
let vm_data vm = vm.data
let vm_state vm = vm.state
let vm_read vm = vm.output
let vm_write x vm = {vm with input=x}
let vm_incr n vm = {vm with pc=vm.pc+n}

let ld x md vm = match md with
  | Pos -> get x vm.data
  | Imm -> x
  | Rel -> get (vm.rel+x) vm.data

let store a pos md vm = match md with
  | Pos -> {vm with data=set a pos vm.data}
  | Imm -> failwith "invalid store in immediate mode"
  | Rel -> {vm with data=set a (vm.rel+pos) vm.data}

let rec step ({pc; data; state} as vm) =
  let {op; modes=m1,m2,m3} = get pc data |> decode in
  let (a,b,c) = (get (pc+1) data), (get (pc+2) data), (get (pc+3) data) in
  let (x,y) = (ld a m1 vm), (ld b m2 vm) in
  if state = Input then 
    let vm = store vm.input a m1 vm in
    step {vm with pc=pc+2; state=Running}
  else match op with
    | Add -> store (x+y) c m3 vm |> vm_incr 4
    | Mul -> store (x*y) c m3 vm |> vm_incr 4
    | Read -> {vm with state=Input}
    | Write -> {vm with pc=pc+2; output=x; state=Output}
    | JmpIf -> {vm with pc=if x<>0 then y else pc+3}
    | JmpNot -> {vm with pc=if x=0 then y else pc+3}
    | Lt -> store (if x<y then 1 else 0) c m3 vm |> vm_incr 4
    | Eq -> store (if x=y then 1 else 0) c m3 vm |> vm_incr 4
    | AdjRel -> {vm with pc=pc+2; rel=vm.rel+x}
    | Halt -> {vm with pc=pc+1; state=Halted}

let rec run vm = match step vm with
  | {state=Running} as vm' -> run vm'
  | vm' -> vm'
