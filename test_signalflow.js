import signalflow from 'k6/x/signalflow';

export default function () {
    const program = "data('demo.trans.latency').publish()";
    const now = Date.now();
    const oneHourAgo = now - (6 * 10 * 1000);
    const resolution = 60000;
    
    const signalFlow = signalflow.newSignalFlow('Dj-KsDFzjcOhhEwpe1t6jg', 'rc0');
    let computation = signalFlow.execute(program,  now, oneHourAgo, resolution);

    console.log("Computations")
    //console.log(computation)

    try {
        console.log("Collecting data")
        // const msgs = computation.collect();
        // console.log(msgs)
    } finally {
        if (computation) {
            computation.close();
        }
        signalFlow.close();
    }

}