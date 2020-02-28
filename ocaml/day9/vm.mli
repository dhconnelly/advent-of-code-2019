type data
val read : Scanf.Scanning.in_channel -> data
val set  : int -> int -> data -> data
val get  : int -> data -> int

type state = Running | Halted | Input | Output

type vm
val vm_new   : data -> vm
val vm_data  : vm -> data
val vm_state : vm -> state
val vm_write : int -> vm -> vm
val vm_read  : vm -> int
val run      : vm -> vm
