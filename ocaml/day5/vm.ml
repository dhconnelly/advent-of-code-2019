let read_program ic =
  let next_instr () =
    try Some (Scanf.bscanf ic "%d%c" (fun i _ -> i))
    with End_of_file -> None
  in let rec read_acc acc = match next_instr () with
  | None -> acc
  | Some instr -> read_acc (instr::acc)
  in read_acc [] |> List.rev |> Array.of_list
