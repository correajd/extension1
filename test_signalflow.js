import signalflow from 'k6/x/signalflow';

export default async function () {
    const realm = __ENV.O11Y_REALM || CONFIG.DEFAULT_REALM;
    const token = __ENV.O11Y_TOKEN;

    const program = "data('demo.trans.latency').publish()";
    const now = Date.now();
    const oneHourAgo = now - (60 * 30 * 1000);
    const resolution = `60000`; // 1 minute in milliseconds
    
    let computation = null;
    let signalFlow = null;
    
    try {
        console.log('Creating SignalFlow client...');
        signalFlow = signalflow.newSignalFlow(token, realm);
        
        console.log(`Executing SignalFlow program from ${new Date(oneHourAgo).toISOString()} to ${new Date(now).toISOString()}`);
        console.log(`Program: ${program}`);
        
        computation = signalFlow.execute(program, oneHourAgo, now, resolution);
        
        if (!computation) {
            throw new Error('Failed to create computation');
        }
        
        console.log('Computation started, waiting for data...');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        console.log('Collecting data...');
        const result = computation.collect();
        
        if (result.error) {
            console.error(`Error during data collection: ${result.error}`);
            return;
        }
        
        const data = result.data || {};
        const tsIDs = Object.keys(data);
        
        if (tsIDs.length === 0) {
            console.log('No data points received');
            return;
        }
        
        console.log(`Collected data for ${tsIDs.length} time series`);
        
        tsIDs.forEach(tsID => {
            const points = data[tsID] || [];
            console.log(`\nTime Series: ${tsID}`);
            console.log(`  Data points: ${points.length}`);
            
            if (points.length > 0) {
                const firstPoint = points[0];
                const lastPoint = points[points.length - 1];
                console.log(`  Time range: ${new Date(firstPoint.timestamp).toISOString()} to ${new Date(lastPoint.timestamp).toISOString()}`);
                console.log(`  First value: ${firstPoint.value}, Last value: ${lastPoint.value}`);
            }

            points.forEach(point => {
                console.log(`    ${new Date(point.timestamp).toISOString()}: ${point.value}`);
            });

        });
        
    } catch (error) {
        console.error('Error in SignalFlow test:', error.message);
    } finally {
        if (computation) {
            console.log('Closing computation...');
            try { computation.close(); } catch (e) {}
        }
        if (signalFlow) {
            console.log('Closing SignalFlow client...');
            try {
                signalFlow.close();
            } catch (closeError) {
                console.error('Error closing SignalFlow client:', closeError.message);
            }
        }
    }
}
