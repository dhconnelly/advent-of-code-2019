open Printf
open Hashtbl

type data = (int, int) Hashtbl.t

let get pos data = match find_opt data pos with
| None -> 0
| Some x -> x

let set x pos data = replace data pos x

let copy = Hashtbl.copy

let read ic =
  let next_instr () =
    try Some (Scanf.bscanf ic "%d%c" (fun i _ -> i))
    with End_of_file -> None
  in let rec read_acc acc = match next_instr () with
  | None -> acc
  | Some instr -> read_acc (instr::acc)
  in read_acc [] |> List.rev
  |> List.mapi (fun i x -> i, x)
  |> List.to_seq |> of_seq

type opcode = Add | Mul | Read | Write
            | JmpIf | JmpNot | Lt | Eq | Halt
type mode = Pos | Imm
type state = Running | Halted | Input | Output

type instruction = {
  op: opcode;
  modes: mode*mode*mode;
}

let op_of x = x mod 100

let parse_mode = function
  | 0 -> Pos
  | 1 -> Imm
  | n -> failwith (sprintf "invalid mode: %d" n)

let modes_of x =
  let x = x / 100 in
  let (m1,m2,m3) = (x mod 10), (x/10 mod 10), (x/100 mod 10) in
 (parse_mode m1), (parse_mode m2), (parse_mode m3)

let decode x =
  let modes = modes_of x in
  match op_of x with
  | 1  -> {op=Add; modes}
  | 2  -> {op=Mul; modes}
  | 3  -> {op=Read; modes}
  | 4  -> {op=Write; modes}
  | 5  -> {op=JmpIf; modes}
  | 6  -> {op=JmpNot; modes}
  | 7  -> {op=Lt; modes}
  | 8  -> {op=Eq; modes}
  | 99 -> {op=Halt; modes}
  | o -> failwith (sprintf "invalid opcode: %d" o)

type vm = {
  pc: int;
  data: data;
  state: state;
  input: int;
  output: int;
}

let vm_new data = {pc=0; data=copy data; input=0; output=0; state=Running}
let vm_data vm = vm.data
let vm_state vm = vm.state
let vm_read vm = vm.output
let vm_write x vm = {vm with input=x}

let ld x md data = match md with
  | Pos -> get x data
  | Imm -> x

let store y x md vm = match md with
  | Pos -> set y x vm.data
  | Imm -> failwith "invalid store in immediate mode"

let rec step ({pc; data; state} as vm) =
  let {op; modes=(m1, m2, m3)} = get pc data |> decode in
  let (p1, p2, p3) = (get (pc+1) data), (get (pc+2) data), (get (pc+3) data) in
  let (a, b) = (ld p1 m1 data), (ld p2 m2 data) in
  if state = Input then (
    store vm.input p1 m1 vm;
    step {vm with pc=pc+2; state=Running}
  ) else match op with
    | Add -> store (a + b) p3 m3 vm; {vm with pc=pc+4}
    | Mul -> store (a * b) p3 m3 vm; {vm with pc=pc+4}
    | Read -> {vm with state=Input}
    | Write -> {vm with pc=pc+2; output=a; state=Output}
    | JmpIf -> {vm with pc=if a<>0 then b else pc+3}
    | JmpNot -> {vm with pc=if a=0 then b else pc+3}
    | Lt -> store (if a < b then 1 else 0) p3 m3 vm; {vm with pc=pc+4}
    | Eq -> store (if a = b then 1 else 0) p3 m3 vm; {vm with pc=pc+4}
    | Halt -> {vm with pc=pc+1; state=Halted}

let rec run vm =
  let vm = step vm in match vm with
  | {state=Running} -> run vm
  | _ -> vm
