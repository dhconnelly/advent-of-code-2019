include Vm

let main prog arg = 
  let vm = vm_new prog in
  let rec loop vm =
    let vm' = run vm in match vm_state vm' with
    | Running -> loop vm'
    | Halted -> ()
    | Input -> vm_write arg vm' |> loop
    | Output -> vm_read vm' |> Printf.printf "%d\n"; loop vm'
  in loop vm
  
let () =
  let path = Sys.argv.(1) in
  let prog = Scanf.Scanning.open_in path |> read in
  main prog 1;
  main prog 2
