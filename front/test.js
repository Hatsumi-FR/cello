const constraints = {
    audio: true,
};

navigator.mediaDevices.getUserMedia(constraints)
    .then(function(stream) {
        const context = new AudioContext();
        const audioInput = context.createMediaStreamSource(stream);
        const bufferSize = 2048;
        const scriptNode = context.createScriptProcessor(bufferSize, 1, 1);

        let websocket = new WebSocket('ws://localhost:8080/ws');

        websocket.onopen = function(event) {
            console.log('success connection');
        };

        websocket.onerror = function(event) {
            console.error('failed to connect ws:', event);
        };

        scriptNode.onaudioprocess = function(audioProcessingEvent) {
            const inputBuffer = audioProcessingEvent.inputBuffer;
            const inputData = inputBuffer.getChannelData(0);

            websocket.send(inputData.buffer);
        };

        audioInput.connect(scriptNode);
        scriptNode.connect(context.destination);
    })
    .catch(function(err) {
        console.error('failed to access mic:', err);
    });
