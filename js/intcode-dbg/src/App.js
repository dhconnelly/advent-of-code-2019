import React from "react";
import "./App.css";
import prog from "./prog.js";
import { VM, State } from "intcode";

function ControlPanel(props) {
    return (
        <div>
            <button onClick={props.onStop}>Stop</button>
            <button onClick={props.onStep}>Step</button>
            <button onClick={props.onRun}>Run</button>
        </div>
    );
}

function StateView(props) {
    return (
        <table>
            <thead></thead>
            <tbody>
                <tr>
                    <td>pc</td>
                    <td>{props.pc}</td>
                </tr>
                <tr>
                    <td>state</td>
                    <td>{props.state}</td>
                </tr>
                <tr>
                    <td>sp</td>
                    <td>{props.sp}</td>
                </tr>
                <tr>
                    <td>input</td>
                    <td>{props.input}</td>
                </tr>
                <tr>
                    <td>output</td>
                    <td>{props.output}</td>
                </tr>
            </tbody>
        </table>
    );
}

function InstructionRow(props) {
    return (
        <tr>
            <td>{props.addr}</td>
            <td>{props.opcode}</td>
            <td>{props.arg1}</td>
            <td>{props.arg2}</td>
            <td>{props.arg3}</td>
        </tr>
    );
}

const INSTRUCTIONS_BEFORE = 5;
const INSTRUCTIONS_AFTER = 5;

function InstructionView(props) {
    return (
        <table>
            <thead>
                <tr>
                    <td>pc</td>
                    <td>opcode</td>
                    <td>arg1</td>
                    <td>arg2</td>
                    <td>arg3</td>
                </tr>
            </thead>
            <tbody>
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
                <InstructionRow />
            </tbody>
        </table>
    );
}

function MemoryRow(props) {
    return (
        <tr>
            <td>{props.addr}</td>
            <td>{props.val}</td>
        </tr>
    );
}

const ROWS_BEFORE = 5;
const ROWS_AFTER = 5;

function MemoryView(props) {
    let rows = [];
    for (let i = -ROWS_BEFORE; i <= ROWS_AFTER; i++) {
        let addr = props.addr + i;
        if (addr < 0) {
            rows.push(<MemoryRow key={i} addr="-" val="-" />);
        } else {
            rows.push(
                <MemoryRow key={i} addr={addr} val={props.memory[addr]} />
            );
        }
    }
    return (
        <table>
            <thead>
                <tr>
                    <td>addr</td>
                    <td>val</td>
                </tr>
            </thead>
            <tbody>{rows}</tbody>
        </table>
    );
}

const TICK_INTERVAL_MS = 100;

class App extends React.Component {
    constructor(props) {
        super(props);
        this.handleVMWrite = this.handleVMWrite.bind(this);
        this.handleStep = this.handleStep.bind(this);
        this.handleStop = this.handleStop.bind(this);
        this.handleRun = this.handleRun.bind(this);
        this.vm = new VM(prog, { onWrite: this.handleVMWrite });
        this.timer = null;
        this.state = {
            pc: this.vm.pc,
            sp: this.vm.sp,
            input: this.vm.input,
            output: this.vm.output,
            memory: this.vm.mem.slice(),
            state: this.vm.state,
            writeLoc: 0,
        };
    }

    handleVMWrite(pos, val) {
        this.setState((prev) => {
            // TODO: use an immutable array for this
            let mem = prev.memory.slice();
            mem[pos] = val;
            return { memory: mem, writeLoc: pos };
        });
    }

    handleStop() {
        clearInterval(this.timer);
    }

    componentWillUnmount() {
        if (this.timer) clearInterval(this.timer);
    }

    handleStep() {
        switch (this.vm.state) {
            // TODO: handle this interactively
            case State.READ:
                this.vm.write(1);
                break;
            case State.WRITE:
                this.vm.read();
                break;
            case State.HALT:
                break;
            case State.RUN:
                this.vm.step();
                this.setState({
                    pc: this.vm.pc,
                    sp: this.vm.sp,
                    input: this.vm.input,
                    output: this.vm.output,
                    state: this.vm.state,
                });
                break;
            default:
                throw new Error("invalid intcode vm state:" + this.vm.state);
        }
    }

    handleRun() {
        this.timer = setInterval(this.handleStep, TICK_INTERVAL_MS);
    }

    render() {
        return (
            <div>
                <ControlPanel
                    onStop={this.handleStop}
                    onStep={this.handleStep}
                    onRun={this.handleRun}
                />
                <StateView
                    pc={this.state.pc}
                    sp={this.state.sp}
                    input={this.state.input}
                    output={this.state.output}
                    state={this.state.state}
                />
                <InstructionView
                    memory={this.state.memory}
                    pc={this.state.pc}
                />
                <MemoryView
                    memory={this.state.memory}
                    addr={this.state.writeLoc}
                />
            </div>
        );
    }
}

export default App;
