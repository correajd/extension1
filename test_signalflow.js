import signalflow from 'k6/x/signalflow';

export default function () {
    const program = "data('demo.trans.latency').publish()";
    const now = Date.now();
    const oneHourAgo = now - (6 * 10 * 1000);
    const resolution = 60000;
    
    const signalFlow = signalflow.newSignalFlow('Dj-KsDFzjcOhhEwpe1t6jg', 'rc0')
    let computation = signalFlow.execute(program,  now, oneHourAgo, resolution);

    console.log("Computations")
    console.log(computation)

    try {
        while (true) {
            const msg = computation.next();
            console.log(`Message: ${JSON.stringify(msg)}`);
            if (!msg) break; // No more messages

            // Handle the message based on its type
            if (msg.type === 'data') {
                // Process data message
                console.log('Data message:', msg);
            } else if (msg.type === 'error') {
                console.error('Error message:', msg);
                break;
            }
            // Add other message type handlers as needed
        }
    } finally {
        computation.close();
        client.close();
    }

}