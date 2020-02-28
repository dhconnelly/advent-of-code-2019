type data
val read : Scanf.Scanning.in_channel -> data
val copy : data -> data
val set  : int -> int -> data -> data
val get  : int -> data -> int

type vm
val vm_new   : data -> vm
val vm_data  : vm -> data
val run      : vm -> vm
